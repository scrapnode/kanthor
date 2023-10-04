package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
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

		for i, event := range events {
			if event.Is(streaming.TopicMsg) {
				msg, err := eventToMessage(i, event)
				if err != nil {
					service.logger.Errorw(err.Error(), "event_id", event.String())
					// IMPORTANT: sometime under high-preasure, we got nil pointer
					// should retry this event anyway
					errs.Set(event.Id, err)
					continue
				}

				messages = append(messages, *msg)
				entity2eventMaps[msg.Id] = event.Id
				continue
			}

			if event.Is(streaming.TopicReq) {
				req, err := eventToRequest(i, event)
				if err != nil {
					service.logger.Errorw(err.Error(), "event_id", event.String())
					// IMPORTANT: sometime under high-preasure, we got nil pointer
					// should retry this event anyway
					errs.Set(event.Id, err)
					continue
				}

				requests = append(requests, *req)
				entity2eventMaps[req.Id] = event.Id
				continue
			}

			if event.Is(streaming.TopicRes) {
				res, err := eventToResponse(i, event)
				if err != nil {
					service.logger.Errorw(err.Error(), "event", event.String())
					// IMPORTANT: sometime under high-preasure, we got nil pointer
					// should retry this event anyway
					errs.Set(event.Id, err)
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
				service.logger.Errorw("encoutered errors", "error_count", errs.Count(), "error_sample", errs.Sample().Error())
			}
			service.logger.Infow("save entities", "entity_count", len(events), "save_count", len(events)-errs.Count())

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

func eventToMessage(i int, event *streaming.Event) (*entities.Message, error) {
	msg, err := transformation.EventToMessage(event)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("events(message).[%d]", i)
	err = validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", msg.Id, "msg_"),
		validator.NumberGreaterThan(prefix+".timestamp", msg.Timestamp, 0),
		validator.StringRequired(prefix+".tier", msg.Tier),
		validator.StringStartsWith(prefix+".app_id", msg.AppId, "app_"),
		validator.StringRequired(prefix+".type", msg.Type),
		validator.SliceRequired(prefix+".body", msg.Body),
	)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func eventToRequest(i int, event *streaming.Event) (*entities.Request, error) {
	req, err := transformation.EventToRequest(event)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("events(request).[%d]", i)
	err = validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", req.Id, "req_"),
		validator.NumberGreaterThan(prefix+".timestamp", req.Timestamp, 0),
		validator.StringStartsWith(prefix+".msg_id", req.MsgId, "msg_"),
		validator.StringStartsWith(prefix+".ep_id", req.EpId, "ep_"),
		validator.StringRequired(prefix+".tier", req.Tier),
		validator.StringStartsWith(prefix+".app_id", req.AppId, "app_"),
		validator.StringRequired(prefix+".type", req.Type),
		validator.SliceRequired(prefix+".body", req.Body),
		validator.StringRequired(prefix+".uri", req.Uri),
		validator.StringRequired(prefix+".method", req.Method),
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func eventToResponse(i int, event *streaming.Event) (*entities.Response, error) {
	res, err := transformation.EventToResponse(event)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("events(response).[%d]", i)
	err = validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", res.Id, "res_"),
		validator.NumberGreaterThan(prefix+".timestamp", res.Timestamp, 0),
		validator.StringStartsWith(prefix+".msg_id", res.MsgId, "msg_"),
		validator.StringStartsWith(prefix+".ep_id", res.EpId, "ep_"),
		validator.StringStartsWith(prefix+".req_id", res.ReqId, "req_"),
		validator.StringRequired(prefix+".tier", res.Tier),
		validator.StringStartsWith(prefix+".app_id", res.AppId, "app_"),
		validator.StringRequired(prefix+".type", res.Type),
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
