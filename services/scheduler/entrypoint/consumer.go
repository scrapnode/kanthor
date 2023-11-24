package entrypoint

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/scheduler/usecase"
)

func NewConsumer(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		messages := map[string]*entities.Message{}
		for _, event := range events {
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

			messages[event.Id] = message
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(service.conf.Request.Schedule.Timeout))
		defer cancel()

		in := &usecase.RequestScheduleIn{
			Messages: messages,
		}
		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Request().Schedule(ctx, in)
		if err != nil {
			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		service.logger.Infow("scheduled requests for messages", "event_count", len(events), "ok_count", len(out.Success), "ko_count", len(out.Error))

		return out.Error
	}
}
