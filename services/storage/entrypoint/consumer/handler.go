package consumer

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/domain/constants"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/transformation"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/storage/usecase"
)

func Handler(service *storage) streaming.SubHandler {
	var counter atomic.Uint64

	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		in := &usecase.WarehousePutIn{
			Size:      service.conf.Warehouse.Put.Size,
			Messages:  map[string]*entities.Message{},
			Requests:  map[string]*entities.Request{},
			Responses: map[string]*entities.Response{},
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

				if err := usecase.ValidateWarehousePutInMessage(prefix, message); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate message", "event", event.String(), "message", message.String(), "err", err.Error())
					continue
				}
				in.Messages[id] = message
				continue
			}

			if event.Is(project.Namespace(), constants.TopicRequest) {
				request, err := transformation.EventToRequest(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform request to event", "event", event.String(), "err", err.Error())
					continue
				}

				if err := usecase.ValidateWarehousePutInRequest(prefix, request); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate request", "event", event.String(), "request", request.String(), "err", err.Error())
					continue
				}
				in.Requests[id] = request
				continue
			}

			if event.Is(project.Namespace(), constants.TopicResponse) {
				response, err := transformation.EventToResponse(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("could not transform response to event", "event", event.String(), "err", err.Error())
					continue
				}

				if err := usecase.ValidateWarehousePutInResponse(prefix, response); err != nil {
					// un-recoverable error
					service.logger.Errorw("could not validate response", "event", event.String(), "response", response.String(), "err", err.Error())
					continue
				}
				in.Responses[id] = response
				continue
			}

			err := fmt.Errorf("unrecognized event %s", event.Id)
			service.logger.Warnw(err.Error(), "event", event.String())
		}

		log.Printf("counter -------------------------------------------> %v", counter.Add(uint64(len(in.Requests))))

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(service.conf.Warehouse.Put.Timeout))
		defer cancel()

		// we alreay validated messages, request and response, don't need to validate again
		out, err := service.uc.Warehouse().Put(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to store entities", "error", err.Error())

			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		if len(out.Error) > 0 {
			for ref, err := range out.Error {
				service.logger.Errorw("unable to store entities", "ref", ref, "error", err.Error())
			}
		}

		service.logger.Infow("put entities", "ok_count", len(out.Success), "ko_count", len(out.Error))

		return out.Error
	}
}
