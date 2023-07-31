package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/dispatcher/transformation"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func Consumer(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		service.meter.Count("dispatcher_send_request_total", 1)

		service.logger.Debugw("received event", "event_id", event.Id)
		req, err := transformation.EventToRequest(event)
		if err != nil {
			service.meter.Count("dispatcher_consume_event_error", 1, metric.Label("action", "transformation"))
			service.logger.Error(err)
			return nil
		}

		request := &usecase.ForwarderSendReq{Request: *req}
		response, err := service.uc.Forwarder().Send(context.Background(), request)
		if err != nil {
			service.meter.Count("dispatcher_send_request_error", 1)
			service.logger.Error(err)
			return nil
		}
		// custom handle for error
		if response.Response.Error != "" {
			service.meter.Count("dispatcher_send_request_error", 1)
		}

		service.logger.Debugw("received response", "response_id", response.Response.Id, "response_status", response.Response.Status)
		return nil
	}
}
