package scheduler

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/services/scheduler/transformation"
	"github.com/sourcegraph/conc"
)

func Consumer(service *scheduler) streaming.SubHandler {
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
				msg, err := transformation.EventToMessage(event)
				if err != nil {
					service.metrics.Count(ctx, "scheduler_transform_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse message from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					return
				}

				ucreq := transformation.MsgToArrangeReq(msg)
				if err := ucreq.Validate(); err != nil {
					service.metrics.Count(ctx, "scheduler_arrange_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String(), "msg", msg.String())
					// unable to construct usecase request from message is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					return
				}

				response, err := service.uc.Request().Arrange(ctx, ucreq)
				if err != nil {
					service.metrics.Count(ctx, "scheduler_arrange_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "msg_id", msg.Id)
					errs.Set(event.Id, err)
					return
				}
				service.metrics.Count(ctx, "scheduler_arrange_total", int64(len(response.Entities)))

				// @TODO: use dead-letter
				if len(response.FailKeys) > 0 {
					service.metrics.Count(ctx, "scheduler_arrange_error", int64(len(response.FailKeys)))
					service.logger.Errorw("got some errors", "fail_keys", response.FailKeys)
				}

				service.logger.Debugw("scheduled requests", "success_count", len(response.SuccessKeys))
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
				service.metrics.Count(ctx, "scheduler_send_error", int64(errs.Count()))
				service.logger.Errorw("encoutered errors", "error_count", errs.Count(), "error_sample", errs.Sample().Error())
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
			service.metrics.Count(ctx, "scheduler_send_timeout_error", 1)
			service.metrics.Count(ctx, "scheduler_send_error", int64(errs.Count()))
			service.logger.Errorw("encoutered errors", "error_count", errs.Count())
			return errs.Data()
		}
	}
}
