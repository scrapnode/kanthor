package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func Consumer(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		service.meter.Counter("dispatcher_consume_event_total", 1)
		service.meter.Counter("dispatcher_send_request_total", 1)

		service.logger.Debugw("received event", "event_id", event.Id)
		req, err := transformEventToRequest(event)
		if err != nil {
			service.meter.Counter("dispatcher_consume_event_error", 1, metric.UseLabel("action", "transform"))
			service.logger.Error(err)
			return nil
		}

		request := &usecase.SendRequestsReq{Request: *req}
		response, err := service.uc.SendRequest(context.TODO(), request)
		if err != nil {
			service.meter.Counter("dispatcher_consume_event_error", 1)
			service.logger.Error(err)
			return nil
		}

		service.logger.Debugw("received response", "response_id", response.Response.Id, "response_status", response.Response.Status)
		return nil
	}
}

func transformEventToRequest(event *streaming.Event) (*entities.Request, error) {
	var req entities.Request
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}
