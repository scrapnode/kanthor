package trigger

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

func RegisterConsumer(service *executor) streaming.SubHandler {
	return func(events map[string]*streaming.Event) map[string]error {
		ucreq := &usecase.TriggerExecReq{
			Size:         service.conf.Attempt.Trigger.Executor.Size,
			Timeout:      service.conf.Attempt.Trigger.Executor.Timeout,
			AttemptDelay: service.conf.Attempt.Trigger.Executor.AttemptDelay,
			Triggers:     map[string]*entities.AttemptTrigger{},
		}

		for _, event := range events {
			trigger, err := transformation.EventToTrigger(event)
			if err != nil {
				service.logger.Errorw("unable to transform event to attempt trigger", "err", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			ucreq.Triggers[trigger.AppId] = trigger
		}

		ctx := context.Background()
		retruning := map[string]error{}

		ucres, err := service.uc.Trigger().Exec(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to consume attempt trigger", "err", err.Error())
			// basically we will not try to retry an attempt trigger
			// because it could be retry later by cronjob
			return retruning
		}

		if len(ucres.Error) > 0 {
			// basically we will not try to retry an attempt trigger
			// because it could be retry later by cronjob
			for key, err := range ucres.Error {
				service.logger.Errorw("consume attempt trigger got some errors", "key", key, "err", err.Error())
			}
		}

		return retruning
	}
}
