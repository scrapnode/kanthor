package consumer

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/transformation"
	"github.com/scrapnode/kanthor/services/scheduler/usecase"
)

func Handler(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		messages := map[string]*entities.Message{}
		for id, event := range events {
			message, err := transformation.EventToMessage(event)
			if err != nil {
				service.logger.Errorw(err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateRequestScheduleInMessage("message", message); err != nil {
				service.logger.Errorw(err.Error(), "event", event.String(), "message", message.String())
				// got malformed message, should ignore and not retry it
				continue
			}

			messages[id] = message
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(service.conf.Request.Schedule.Timeout))
		defer cancel()

		in := &usecase.RequestScheduleIn{
			Messages: messages,
		}
		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Request().Schedule(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to schedule requests", "error", err.Error())

			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		service.logger.Infow("scheduled requests for messages", "event_count", len(events), "ok_count", len(out.Success), "ko_count", len(out.Error))
		if len(out.Error) > 0 {
			for ref, err := range out.Error {
				service.logger.Errorw("schedule got error", "ref", ref, "error", err.Error())
			}
		}

		return out.Error
	}
}
