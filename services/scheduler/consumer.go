package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/usecases"
)

func Consumer(logger logging.Logger, usecase usecases.Scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		msg, err := transformEventToMessage(event)
		if err != nil {
			logger.Error(err)
			return nil
		}

		req := &usecases.ArrangeRequestsReq{Message: msg}
		res, err := usecase.ArrangeRequests(context.TODO(), req)
		if err != nil {
			logger.Error(err)
			return nil
		}

		// @TODO: use deadletter
		if len(res.FailKeys) > 0 {
			logger.Errorw("got some errors", "fail_keys", res.FailKeys)
		}

		logger.Debugw("scheduled requested", "request_count", len(res.SuccessKeys))
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
