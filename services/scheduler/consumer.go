package scheduler

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

func NewConsumer(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []*streaming.Event) map[string]error {
		// create a map of events & messages so we can generate error map later
		maps := map[string]string{}

		messages := []entities.Message{}
		for _, event := range events {
			message, err := transformation.EventToMessage(event)
			if err != nil {
				service.logger.Errorw(err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateRequestScheduleMessaeg("message", *message); err != nil {
				service.logger.Errorw(err.Error(), "event", event.String(), "message", message.String())
				// got malformed message, should ignore and not retry it
				continue
			}

			maps[message.Id] = event.Id
			messages = append(messages, *message)
		}

		retruning := map[string]error{}
		ctx := context.Background()

		ucreq := &usecase.RequestScheduleReq{
			ChunkTimeout: service.conf.Scheduler.Request.Schedule.ChunkTimeout,
			ChunkSize:    service.conf.Scheduler.Request.Schedule.ChunkSize,
			Messages:     messages,
		}
		// we alreay validated messages of request, don't need to validate again
		ucres, err := service.uc.Request().Schedule(ctx, ucreq)
		if err != nil {
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		// ucres.Error contain a map of message id and error
		// should convert it to a map of event id and error so our streaming service can retry it
		if len(ucres.Error) > 0 {
			for msgId, err := range ucres.Error {
				eventId := maps[msgId]
				retruning[eventId] = err
			}
		}

		return retruning
	}
}
