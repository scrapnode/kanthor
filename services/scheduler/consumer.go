package scheduler

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/scheduler/transformation"
	"github.com/sourcegraph/conc"
)

func Consumer(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []streaming.Event) map[string]error {
		results := map[string]error{}

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
					service.logger.Error(err)
					// next retry will be failed because of same error of transformation
					return
				}

				request := transformation.MsgToArrangeReq(msg)
				if err := service.validator.Struct(request); err != nil {
					service.metrics.Count(ctx, "dispatcher_arrange_error", 1)
					service.logger.Errorw(err.Error(), "data", event.String())
					results[event.Id] = err
					return
				}

				response, err := service.uc.Request().Arrange(ctx, request)
				if err != nil {
					service.metrics.Count(ctx, "dispatcher_arrange_error", 1)
					service.logger.Errorw(err.Error(), "data", event.String())
					results[event.Id] = err
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
