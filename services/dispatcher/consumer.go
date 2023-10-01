package dispatcher

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/services/dispatcher/transformation"
	"github.com/sourcegraph/conc"
)

func Consumer(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []*streaming.Event) map[string]error {
		errs := &ds.SafeMap[error]{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

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

				ucreq := transformation.ReqToSendReq(req)
				if err := ucreq.Validate(); err != nil {
					service.metrics.Count(ctx, "dispatcher_send_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String(), "req", req.String())
					// unable to construct usecase request from request is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
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
