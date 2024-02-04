package selector

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func Handler(service *selector) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(ctx context.Context, events map[string]*streaming.Event) map[string]error {
		timeout := time.Millisecond * time.Duration(service.conf.Consumer.Timeout)
		timeoutctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		in := &usecase.RetrySelectIn{
			BatchSize: service.conf.Selector.BatchSize,
			Counter:   service.conf.Selector.Counter,
			Triggers:  make(map[string]*entities.AttemptTrigger),
		}

		for id, event := range events {
			trigger, err := transformation.EventToAttemptTrigger(event)
			if err != nil {
				service.logger.Errorw("ATTEMPT.ENTRYPOINT.SELECTOR.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
				// unable to parse attempt from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateRetrySelectAttemptTrigger("trigger", trigger); err != nil {
				service.logger.Errorw("ATTEMPT.ENTRYPOINT.SELECTOR.HANDLER.ATTEMPT_TASK_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "attempt_trigger", trigger.String())
				// got malformed attempt, should ignore and not retry it
				continue
			}

			in.Triggers[id] = trigger
		}

		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Retry().Select(timeoutctx, in)
		if err != nil {
			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		return out.Error
	}
}
