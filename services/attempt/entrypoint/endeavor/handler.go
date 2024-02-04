package endeavor

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func Handler(service *endeavor) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(ctx context.Context, events map[string]*streaming.Event) map[string]error {
		timeout := time.Millisecond * time.Duration(service.conf.Consumer.Timeout)
		timeoutctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		in := &usecase.RetryEndeavorIn{
			Concurrency: service.conf.Endeavor.Concurrency,
			Attempts:    make(map[string]*entities.Attempt),
		}

		for id, event := range events {
			attempt, err := transformation.EventToAttempt(event)
			if err != nil {
				service.logger.Errorw("ATTEMPT.ENTRYPOINT.ENDEAVOR.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
				// unable to parse attempt from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateRetryEndeavorInAttempt("attempt", attempt); err != nil {
				service.logger.Errorw("ATTEMPT.ENTRYPOINT.ENDEAVOR.HANDLER.ATTEMPT_TASK_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "attempt", attempt.String())
				// got malformed attempt, should ignore and not retry it
				continue
			}

			in.Attempts[id] = attempt
		}

		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Retry().Endeavor(timeoutctx, in)
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
