package consumer

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/services/scheduler/usecase"
)

func Handler(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(ctx context.Context, events map[string]*streaming.Event) map[string]error {
		// timeout
		timeout := time.Millisecond * time.Duration(service.conf.Request.Schedule.Timeout)
		timeoutctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		in := &usecase.RequestScheduleIn{
			Messages: make(map[string]*entities.Message),
		}
		for id, event := range events {
			message, err := transformation.EventToMessage(event)
			if err != nil {
				service.logger.Errorw("SCHEDULER.ENTRYPOINT.CONSUMER.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateRequestScheduleInMessage("message", message); err != nil {
				service.logger.Errorw("SCHEDULER.ENTRYPOINT.CONSUMER.HANDLER.MESSAGE_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "message", message.String())
				// got malformed message, should ignore and not retry it
				continue
			}

			in.Messages[id] = message
		}

		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Request().Schedule(timeoutctx, in)
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
