package consumer

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/services/dispatcher/usecase"
)

func Handler(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		requests := map[string]*entities.Request{}
		for id, event := range events {
			request, err := transformation.EventToRequest(event)
			if err != nil {
				service.logger.Errorw(err.Error(), "event", event.String())
				// unable to parse request from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateForwarderSendInRequest("request", request); err != nil {
				service.logger.Errorw(err.Error(), "event", event.String(), "request", request.String())
				// got malformed request, should ignore and not retry it
				continue
			}

			requests[id] = request
		}

		ctx := context.Background()
		in := &usecase.ForwarderSendIn{
			Concurrency: service.conf.Forwarder.Send.Concurrency,
			Requests:    requests,
		}
		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Forwarder().Send(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to dispatch requests", "error", err.Error())

			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for refId := range in.Requests {
				retruning[refId] = err
			}
			return retruning
		}

		service.logger.Infow("dispatched requests", "event_count", len(events), "ok_count", len(out.Success), "ko_count", len(out.Error))
		if len(out.Error) > 0 {
			for ref, err := range out.Error {
				service.logger.Errorw("dispatch got error", "ref", ref, "error", err.Error())
			}
		}

		return out.Error
	}
}
