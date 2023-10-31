package entrypoint

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/dispatcher/usecase"
)

func NewConsumer(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		requests := map[string]*entities.Request{}
		for _, event := range events {
			request, err := transformation.EventToRequest(event)
			if err != nil {
				service.logger.Errorw(err.Error(), "event", event.String())
				// unable to parse request from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateForwarderSendReqRequest("request", request); err != nil {
				service.logger.Errorw(err.Error(), "event", event.String(), "request", request.String())
				// got malformed request, should ignore and not retry it
				continue
			}

			requests[event.Id] = request
		}

		ctx := context.Background()
		ucreq := &usecase.ForwarderSendReq{
			Concurrency: service.conf.Forwarder.Send.Concurrency,
			Requests:    requests,
		}
		// we alreay validated messages of request, don't need to validate again
		ucres, err := service.uc.Forwarder().Send(ctx, ucreq)
		if err != nil {
			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for refId := range ucreq.Requests {
				retruning[refId] = err
			}
			return retruning
		}

		return ucres.Error
	}
}