package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func Consumer(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		service.meter.Count("scheduler_arrange_request_total", 1)

		service.logger.Debugw("received event", "event_id", event.Id)
		msg, err := transformEventToMessage(event)
		if err != nil {
			service.meter.Count("scheduler_arrange_request_error", 1, metric.Label("action", "transform"))
			service.logger.Error(err)
			return nil
		}

		request := &usecase.ArrangeRequestsReq{Message: *msg}
		response, err := service.uc.ArrangeRequests(context.TODO(), request)
		if err != nil {
			service.meter.Count("scheduler_arrange_request_error", 1)
			service.logger.Error(err)
			return nil
		}

		service.meter.Count("scheduler_arrange_request_entity_total", int64(len(response.Entities)))
		// @TODO: use deadletter
		if len(response.FailKeys) > 0 {
			service.meter.Count("scheduler_arrange_request_entity_fail_total", int64(len(response.FailKeys)))
			service.logger.Errorw("got some errors", "fail_keys", response.FailKeys)
		}

		service.logger.Debugw("scheduled requested", "success_count", len(response.SuccessKeys))
		return nil
	}
}

func transformEventToMessage(event *streaming.Event) (*entities.Message, error) {
	var msg entities.Message
	if err := msg.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &msg, nil
}
