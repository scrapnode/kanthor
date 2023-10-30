package executor

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func RegisterConsumer(service *executor) streaming.SubHandler {
	return func(events map[string]*streaming.Event) map[string]error {
		ucreq := &usecase.EndeavorExecReq{
			Concurrency: service.conf.Endeavor.Executor.Concurrency,
			Attempts:    map[string]*entities.Attempt{},
		}

		for _, event := range events {
			attempt, err := transformation.EventToAttempt(event)
			if err != nil {
				service.logger.Errorw("unable to transform event to attempt", "err", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			ucreq.Attempts[event.Id] = attempt
		}

		ctx := context.Background()

		ucres, err := service.uc.Endeavor().Exec(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to consume an attempt", "err", err.Error())
			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for refId := range ucreq.Attempts {
				retruning[refId] = err
			}
			return retruning
		}

		return ucres.Error
	}
}
