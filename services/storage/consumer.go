package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/storage"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
)

func Consumer(service *storage) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []*streaming.Event) map[string]error {
		errs := map[string]error{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		service.metrics.Count(ctx, "storage_put_total", int64(len(events)))

		maps := map[string]string{}
		messages := []entities.Message{}
		requests := []entities.Request{}
		responses := []entities.Response{}

		for _, event := range events {
			if event.Is(streaming.TopicMsg) {
				msg, err := transformation.EventToMessage(event)
				if err != nil {
					errs[event.Id] = err
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
					errs[event.Id] = err
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
					errs[event.Id] = err
					service.logger.Error(err)
					// next retry will be failed because of same error of transformation
					continue
				}
				responses = append(responses, *res)
				maps[res.Id] = event.Id
				continue
			}

			err := fmt.Errorf("unrecognized event %s", event.Id)
			errs[event.Id] = err
			service.logger.Warnw(err.Error(), "event", utils.Stringify(event))
		}

		// IMPORTANT: we ignore all validations of storage to trade validity to performance
		var wg conc.WaitGroup

		wg.Go(func() {
			if len(messages) > 0 {
				_, err := service.uc.Message().Put(ctx, &usecase.MessagePutReq{Docs: messages})
				if err != nil {
					for _, msg := range messages {
						eventId := maps[msg.Id]
						errs[eventId] = err
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
						errs[eventId] = err
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
						errs[eventId] = err
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
			if len(errs) > 0 {
				service.metrics.Count(ctx, "storage_put_error", int64(len(errs)))
				service.logger.Errorw("encoutered errors", "error_count", len(errs))
			}
			return errs
		case <-ctx.Done():
			// timeout, all events will be considered as failed
			for _, event := range events {
				if err, ok := errs[event.Id]; ok && err != nil {
					errs[event.Id] = ctx.Err()
				}
			}
			if len(errs) > 0 {
				service.metrics.Count(ctx, "storage_put_timeout_error", 1)
				service.metrics.Count(ctx, "storage_put_error", int64(len(errs)))
				service.logger.Errorw("encoutered errors", "error_count", len(errs))
			}
			return errs
		}
	}
}
