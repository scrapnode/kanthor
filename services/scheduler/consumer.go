package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func Consumer(logger logging.Logger, uc usecase.Scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		logger.Debugw("received event", "event_id", event.Id)
		msg, err := transformEventToMessage(event)
		if err != nil {
			logger.Error(err)
			return nil
		}

		request := &usecase.ArrangeRequestsReq{Message: *msg}
		response, err := uc.ArrangeRequests(context.TODO(), request)
		if err != nil {
			logger.Error(err)
			return nil
		}

		// @TODO: use deadletter
		if len(response.FailKeys) > 0 {
			logger.Errorw("got some errors", "fail_keys", response.FailKeys)
		}

		logger.Debugw("scheduled requested", "request_count", len(response.SuccessKeys))
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
