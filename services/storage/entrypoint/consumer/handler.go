package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/constants"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/storage/usecase"
)

func Handler(service *storage) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(ctx context.Context, events map[string]*streaming.Event) map[string]error {
		timeout := time.Millisecond * time.Duration(service.conf.Warehouse.Put.Timeout)
		timeoutctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		in := &usecase.WarehousePutIn{
			BatchSize: service.conf.Warehouse.Put.BatchSize,
			Messages:  map[string]*entities.Message{},
			Requests:  map[string]*entities.Request{},
			Responses: map[string]*entities.Response{},
			Attempts:  map[string]*entities.Attempt{},
		}

		for id, event := range events {
			prefix := fmt.Sprintf("event.%s", id)

			if project.IsTopic(event.Subject, constants.TopicMessage) {
				message, err := transformation.EventToMessage(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
					continue
				}

				if err := usecase.ValidateWarehousePutInMessage(prefix, message); err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.MESSAGE_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "message", message.String())
					continue
				}
				in.Messages[id] = message
				continue
			}

			if project.IsTopic(event.Subject, constants.TopicRequest) {
				request, err := transformation.EventToRequest(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
					continue
				}

				if err := usecase.ValidateWarehousePutInRequest(prefix, request); err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.REQUEST_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "request", request.String())
					continue
				}
				in.Requests[id] = request
				continue
			}

			if project.IsTopic(event.Subject, constants.TopicResponse) {
				response, err := transformation.EventToResponse(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
					continue
				}

				if err := usecase.ValidateWarehousePutInResponse(prefix, response); err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.RESPONSE_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "response", response.String())
					continue
				}
				in.Responses[id] = response
				continue
			}

			if project.IsTopic(event.Subject, constants.TopicAttempt) {
				attempt, err := transformation.EventToAttempt(event)
				if err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
					continue
				}

				if err := usecase.ValidateWarehousePutInAttempt(prefix, attempt); err != nil {
					// un-recoverable error
					service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.RESPONSE_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "attempt", attempt.String())
					continue
				}
				in.Attempts[id] = attempt
				continue
			}

			service.logger.Errorw("STORAGE.ENTRYPOINT.CONSUMER.HANDLER.EVENT_UNKNOWN_TOPIC.ERROR", "event", event.String())
		}

		// we alreay validated messages, request and response, don't need to validate again
		out, err := service.uc.Warehouse().Put(timeoutctx, in)
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
