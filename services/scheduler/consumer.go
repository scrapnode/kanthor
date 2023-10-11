package scheduler

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/services/scheduler/transformation"
	"github.com/sourcegraph/conc"
)

func NewConsumer(service *scheduler) streaming.SubHandler {
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
				msg, err := transformation.EventToMessage(event)
				if err != nil {
					service.metrics.Count(ctx, "scheduler_transform_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse message from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					return
				}

				ucreq := transformation.MessageToRequestArrangeReq(msg)
				if err := ucreq.Validate(); err != nil {
					service.metrics.Count(ctx, "scheduler_request_arrange_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String(), "msg", msg.String())
					// IMPORTANT: sometime under high-preasure, we got nil pointer
					// should retry this event anyway
					errs.Set(event.Id, err)
					return
				}

				arranged, err := service.uc.Request().Arrange(ctx, ucreq)
				if err != nil {
					service.metrics.Count(ctx, "scheduler_request_arrange_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "msg_id", msg.Id)
					errs.Set(event.Id, err)
					return
				}
				service.metrics.Count(ctx, "scheduler_request_arrange_total", int64(len(arranged.Requests)))

				scheduled, err := service.uc.Request().Schedule(ctx, transformation.RequestToRequestScheduleReq(arranged.Requests))
				if err != nil {
					service.metrics.Count(ctx, "scheduler_request_schedule_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "msg_id", msg.Id)
					errs.Set(event.Id, err)
					return
				}
				service.metrics.Count(ctx, "scheduler_request_schedule_success", int64(len(scheduled.Success)))
				service.metrics.Count(ctx, "scheduler_request_schedule_error", int64(len(scheduled.Error)))

				// @TODO: use dead-letter
				if len(scheduled.Error) > 0 {
					for key, err := range scheduled.Error {
						service.logger.Errorw("schedule error", "key", key, "err", err.Error())
					}
				}

				service.logger.Debugw("scheduled requests", "success_count", len(scheduled.Success))
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
			service.logger.Infow("scheduled requests", "message_count", len(events), "request_count", len(events)-errs.Count())

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
