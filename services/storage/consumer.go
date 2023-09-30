package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/storage"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
)

func Consumer(service *storage) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []*streaming.Event) map[string]error {
		errs := &ds.SafeMap[error]{}

		entity2eventMaps := map[string]string{}
		messages := []entities.Message{}
		requests := []entities.Request{}
		responses := []entities.Response{}

		for _, event := range events {
			if event.Is(streaming.TopicMsg) {
				msg, err := transformation.EventToMessage(event)
				if err != nil {
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse message from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					continue
				}
				messages = append(messages, *msg)
				entity2eventMaps[msg.Id] = event.Id
				continue
			}

			if event.Is(streaming.TopicReq) {
				req, err := transformation.EventToRequest(event)
				if err != nil {
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse request from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					continue
				}
				requests = append(requests, *req)
				entity2eventMaps[req.Id] = event.Id
				continue
			}

			if event.Is(streaming.TopicRes) {
				res, err := transformation.EventToResponse(event)
				if err != nil {
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse response from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					continue
				}
				responses = append(responses, *res)
				entity2eventMaps[res.Id] = event.Id
				continue
			}

			err := fmt.Errorf("unrecognized event %s", event.Id)
			errs.Set(event.Id, err)
			service.logger.Warnw(err.Error(), "event", utils.Stringify(event))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		// IMPORTANT: we ignore all validations of storage to trade validity to performance
		var wg conc.WaitGroup

		if len(messages) > 0 {
			wg.Go(func() {
				_, err := service.uc.Message().Put(ctx, &usecase.MessagePutReq{Docs: messages})
				if err != nil {
					for _, msg := range messages {
						eventId := entity2eventMaps[msg.Id]
						errs.Set(eventId, err)
					}
					return
				}
			})
		}

		if len(requests) > 0 {
			wg.Go(func() {
				_, err := service.uc.Request().Put(ctx, &usecase.RequestPutReq{Docs: requests})
				if err != nil {
					for _, req := range requests {
						eventId := entity2eventMaps[req.Id]
						errs.Set(eventId, err)
					}
					return
				}
			})
		}

		if len(responses) > 0 {
			wg.Go(func() {
				_, err := service.uc.Response().Put(ctx, &usecase.ResponsePutReq{Docs: responses})
				if err != nil {
					for _, res := range responses {
						eventId := entity2eventMaps[res.Id]
						errs.Set(eventId, err)
					}
					return
				}
			})
		}

		c := make(chan bool)
		defer close(c)

		go func() {
			wg.Wait()
			c <- true
		}()

		select {
		case <-c:
			if errs.Count() > 0 {
				service.metrics.Count(ctx, "storage_put_error", int64(errs.Count()))
				service.logger.Errorw("encoutered errors", "error_count", errs.Count())
			}
			return errs.Data()
		case <-ctx.Done():
			// timeout, all events will be considered as failed
			// set non-error event with timeout error
			for _, event := range events {
				if err, ok := errs.Get(event.Id); !ok && err == nil {
					errs.Set(event.Id, ctx.Err())
				}
			}
			service.metrics.Count(ctx, "storage_put_timeout_error", 1)
			service.metrics.Count(ctx, "storage_put_error", int64(errs.Count()))
			service.logger.Errorw("encoutered errors", "error_count", errs.Count())
			return errs.Data()
		}
	}
}
