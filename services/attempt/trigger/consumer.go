package trigger

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/attempt"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

func RegisterConsumer(service *executor) streaming.SubHandler {
	return func(events []*streaming.Event) map[string]error {
		ucreq := &usecase.TriggerExecReq{
			Timeout:       service.conf.Attempt.Trigger.Executor.Timeout,
			RateLimit:     service.conf.Attempt.Trigger.Executor.RateLimit,
			AttemptDelay:  service.conf.Attempt.Trigger.Executor.AttemptDelay,
			Notifications: []entities.AttemptTrigger{},
		}

		for _, event := range events {
			notification, err := transformation.EventToTrigger(event)
			if err != nil {
				service.logger.Errorw("unable to transform event to attempt notification", "err", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			ucreq.Notifications = append(ucreq.Notifications, *notification)
		}

		ctx := context.Background()
		retruning := map[string]error{}

		ucres, err := service.uc.Trigger().Exec(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to consume attempt notifications", "err", err.Error())
			// basically we will not try to retry an attempt notification
			// because it could be retry later by cronjob
			return retruning
		}

		if len(ucres.Error) > 0 {
			// basically we will not try to retry an attempt notification
			// because it could be retry later by cronjob
			for key, err := range ucres.Error {
				service.logger.Errorw("consume attempt notification got some errors", "key", key, "err", err.Error())
			}
		}

		return retruning
	}
}
