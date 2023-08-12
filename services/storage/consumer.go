package storage

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/storage"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
	"time"
)

func Consumer(service *storage) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []streaming.Event) map[string]error {
		results := map[string]error{}

		maps := map[string]string{}
		messages := []entities.Message{}
		requests := []entities.Request{}
		responses := []entities.Response{}

		for _, event := range events {
			if event.Is(streaming.TopicMsg) {
				msg, err := transformation.EventToMessage(event)
				if err != nil {
					service.logger.Error(err)
					// next retry will be failed because of same error of transformation
					continue
				}
				messages = append(messages, *msg)
				maps[msg.Id] = event.Id
				continue
			}

			if event.Is(streaming.TopicReq) {
				req, err := transformation.EventToRequest(event)
				if err != nil {
					service.logger.Error(err)
					// next retry will be failed because of same error of transformation
					continue
				}
				requests = append(requests, *req)
				maps[req.Id] = event.Id
				continue
			}

			if event.Is(streaming.TopicRes) {
				res, err := transformation.EventToResponse(event)
				if err != nil {
					service.logger.Error(err)
					// next retry will be failed because of same error of transformation
					continue
				}
				responses = append(responses, *res)
				maps[res.Id] = event.Id
				continue
			}

			service.logger.Warnw("unrecognized event", "event", utils.Stringify(event))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var wg conc.WaitGroup

		wg.Go(func() {
			if len(messages) > 0 {
				_, err := service.uc.Message().Put(ctx, &usecase.MessagePutReq{Docs: messages})
				if err != nil {
					for _, msg := range messages {
						eventId := maps[msg.Id]
						results[eventId] = err
					}
					return
				}
			}
		})

		wg.Go(func() {
			if len(requests) > 0 {
				_, err := service.uc.Request().Put(ctx, &usecase.RequestPutReq{Docs: requests})
				if err != nil {
					for _, req := range requests {
						eventId := maps[req.Id]
						results[eventId] = err
					}
					return
				}
			}
		})

		wg.Go(func() {
			if len(responses) > 0 {
				_, err := service.uc.Response().Put(ctx, &usecase.ResponsePutReq{Docs: responses})
				if err != nil {
					for _, res := range responses {
						eventId := maps[res.Id]
						results[eventId] = err
					}
					return
				}
			}
		})

		c := make(chan bool)
		go func() {
			defer close(c)
			wg.Wait()
		}()

		select {
		case <-c:
			return results
		case <-ctx.Done():
			for _, event := range events {
				results[event.Id] = ctx.Err()
			}
			return results
		}
	}
}
