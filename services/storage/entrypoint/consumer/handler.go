package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/domain/constants"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/transformation"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/storage/usecase"
)

func Handler(service *storage) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		retruning := map[string]error{}

		// create a map of events & entities so we can generate error map later
		maps := map[string]string{}

		in := &usecase.WarehousePutIn{
			Size:      service.conf.Warehouse.Put.Size,
			Messages:  []entities.Message{},
			Requests:  []entities.Request{},
			Responses: []entities.Response{},
		}

		for id, event := range events {
			prefix := fmt.Sprintf("event[%s]", id)

			if event.Is(project.Namespace(), constants.TopicMessage) {
				message, err := transformation.EventToMessage(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform message to event", "event", event.String(), "err", err.Error())
					continue
				}
				maps[message.Id] = event.Id

				if err := usecase.ValidateWarehousePutInMessage(prefix, message); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate message", "event", event.String(), "message", message.String(), "err", err.Error())
					continue
				}
				in.Messages = append(in.Messages, *message)
				continue
			}

			if event.Is(project.Namespace(), constants.TopicRequest) {
				request, err := transformation.EventToRequest(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform request to event", "event", event.String(), "err", err.Error())
					continue
				}
				maps[request.Id] = event.Id

				if err := usecase.ValidateWarehousePutInRequest(prefix, request); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate request", "event", event.String(), "request", request.String(), "err", err.Error())
					continue
				}
				in.Requests = append(in.Requests, *request)
				continue
			}

			if event.Is(project.Namespace(), constants.TopicResponse) {
				response, err := transformation.EventToResponse(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform response to event", "event", event.String(), "err", err.Error())
					continue
				}
				maps[response.Id] = event.Id

				if err := usecase.ValidateWarehousePutInResponse(prefix, response); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate response", "event", event.String(), "response", response.String(), "err", err.Error())
					continue
				}
				in.Responses = append(in.Responses, *response)
				continue
			}

			err := fmt.Errorf("unrecognized event %s", event.Id)
			retruning[event.Id] = err
			service.logger.Warnw(err.Error(), "event", event.String())
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(service.conf.Warehouse.Put.Timeout))
		defer cancel()

		// we alreay validated messages, request and response, don't need to validate again
		out, err := service.uc.Warehouse().Put(ctx, in)
		if err != nil {
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		// out.Error contain a map of entity id and error
		// should convert it to a map of event id and error so our streaming service can retry it
		if len(out.Error) > 0 {
			for entId, err := range out.Error {
				eventId := maps[entId]
				retruning[eventId] = err
			}
		}

		service.logger.Infow("put entities", "ok_count", len(out.Success), "ko_count", len(out.Error))

		return retruning
	}
}
