package entrypoint

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/storage/usecase"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

func NewConsumer(service *storage) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		retruning := map[string]error{}
		ctx := context.Background()

		// create a map of events & entities so we can generate error map later
		maps := map[string]string{}

		ucreq := &usecase.WarehousePutReq{
			Timeout:   service.conf.Warehouse.Put.Timeout,
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

				if err := usecase.ValidateWarehousePutReqMessage(prefix, message); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate message", "event", event.String(), "message", message.String(), "err", err.Error())
					continue
				}
				ucreq.Messages = append(ucreq.Messages, *message)
			}

			if event.Is(project.Namespace(), constants.TopicRequest) {
				request, err := transformation.EventToRequest(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform request to event", "event", event.String(), "err", err.Error())
					continue
				}
				maps[request.Id] = event.Id

				if err := usecase.ValidateWarehousePutReqRequest(prefix, request); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate request", "event", event.String(), "request", request.String(), "err", err.Error())
					continue
				}
				ucreq.Requests = append(ucreq.Requests, *request)
			}

			if event.Is(project.Namespace(), constants.TopicResponse) {
				response, err := transformation.EventToResponse(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform response to event", "event", event.String(), "err", err.Error())
					continue
				}
				maps[response.Id] = event.Id

				if err := usecase.ValidateWarehousePutReqResponse(prefix, response); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate response", "event", event.String(), "response", response.String(), "err", err.Error())
					continue
				}
				ucreq.Responses = append(ucreq.Responses, *response)
			}

			err := fmt.Errorf("unrecognized event %s", event.Id)
			retruning[event.Id] = err
			service.logger.Warnw(err.Error(), "event", event.String())
		}

		// we alreay validated messages, request and response, don't need to validate again
		ucres, err := service.uc.Warehouse().Put(ctx, ucreq)
		if err != nil {
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		// ucres.Error contain a map of entity id and error
		// should convert it to a map of event id and error so our streaming service can retry it
		if len(ucres.Error) > 0 {
			for entId, err := range ucres.Error {
				eventId := maps[entId]
				retruning[eventId] = err
			}
		}

		return retruning
	}
}
