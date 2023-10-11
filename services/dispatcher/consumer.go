package dispatcher

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
)

func NewConsumer(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []*streaming.Event) map[string]error {
		errs := &ds.SafeMap[error]{}

		// @TODO: remove hardcode timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// @TODO: use pool with WithContext and max goroutine is number of events
		// if err := p.Wait(); err != nil {
		// 	return nil, err
		// }

		var wg conc.WaitGroup
		for _, e := range events {
			event := e
			wg.Go(func() {
				service.logger.Debugw("received event", "event_id", event.Id)
				req, err := transformation.EventToRequest(event)
				if err != nil {
					service.metrics.Count(ctx, "dispatcher_transform_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse request from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					return
				}

				ucreq := &usecase.ForwarderSendReq{
					Request: usecase.ForwarderSendReqRequest{
						Id:       req.Id,
						MsgId:    req.MsgId,
						EpId:     req.EpId,
						Tier:     req.Tier,
						AppId:    req.AppId,
						Type:     req.Type,
						Metadata: req.Metadata,
						Headers:  req.Headers,
						Body:     req.Body,
						Uri:      req.Uri,
						Method:   req.Method,
					},
				}
				if err := ucreq.Validate(); err != nil {
					service.metrics.Count(ctx, "dispatcher_send_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String(), "req", req.String())
					// IMPORTANT: sometime under high-preasure, we got nil pointer
					// should retry this event anyway
					errs.Set(event.Id, err)
					return
				}

				response, err := service.uc.Forwarder().Send(ctx, ucreq)
				if err != nil {
					service.metrics.Count(ctx, "dispatcher_send_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "req_id", req.Id)
					errs.Set(event.Id, err)
					return
				}

				service.logger.Debugw("received response", "res_id", response.Response.Id, "res_status", response.Response.Status)
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
				service.metrics.Count(ctx, "dispatcher_send_error", int64(errs.Count()))
				service.logger.Errorw("encoutered errors", "error_count", errs.Count(), "error_sample", errs.Sample().Error())
			}
			service.logger.Infow("send requests", "request_count", len(events), "response_count", len(events)-errs.Count())

			return errs.Data()
		case <-ctx.Done():
			// timeout, all events will be considered as failed
			// set non-error event with timeout error
			for _, event := range events {
				if err, ok := errs.Get(event.Id); !ok && err == nil {
					errs.Set(event.Id, ctx.Err())
				}
			}
			service.metrics.Count(ctx, "dispatcher_send_timeout_error", 1)
			service.metrics.Count(ctx, "dispatcher_send_error", int64(errs.Count()))
			service.logger.Errorw("encoutered errors", "error_count", errs.Count())
			return errs.Data()
		}
	}
}
