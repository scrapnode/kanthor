package executor

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func Handler(service *executor) streaming.SubHandler {
	return func(events map[string]*streaming.Event) map[string]error {
		in := &usecase.EndeavorExecIn{
			Concurrency: service.conf.Endeavor.Executor.Concurrency,
			Attempts:    map[string]*entities.Attempt{},
		}

		for id, event := range events {
			attempt, err := transformation.EventToAttempt(event)
			if err != nil {
				service.logger.Errorw("unable to transform event to attempt endeavor", "err", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			in.Attempts[id] = attempt
		}

		ctx := context.Background()

		out, err := service.uc.Endeavor().Exec(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to consume an attempt endeavors", "err", err.Error())
			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for refId := range in.Attempts {
				retruning[refId] = err
			}
			return retruning
		}

		service.logger.Infow("consumed attempt endeavors", "event_count", len(events), "ok_count", len(out.Success), "ko_count", len(out.Error))
		if len(out.Error) > 0 {
			for ref, err := range out.Error {
				service.logger.Errorw("endeavor got error", "ref", ref, "error", err.Error())
			}
		}

		return out.Error
	}
}
