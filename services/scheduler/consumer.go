package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/scheduler/transformation"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func Consumer(service *scheduler) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		service.logger.Debugw("received event", "event_id", event.Id)
		msg, err := transformation.EventToMessage(event)
		if err != nil {
			service.logger.Error(err)
			return nil
		}

		request := &usecase.RequestArrangeReq{Message: *msg}
		response, err := service.uc.Request().Arrange(context.TODO(), request)
		if err != nil {
			service.logger.Error(err)
			return nil
		}

		// @TODO: use dead-letter
		if len(response.FailKeys) > 0 {
			service.logger.Errorw("got some errors", "fail_keys", response.FailKeys)
		}

		service.logger.Debugw("scheduled requests", "success_count", len(response.SuccessKeys))
		return nil
	}
}
