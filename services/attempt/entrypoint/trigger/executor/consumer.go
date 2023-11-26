package executor

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/transformation"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func RegisterConsumer(service *executor) streaming.SubHandler {
	return func(events map[string]*streaming.Event) map[string]error {
		in := &usecase.TriggerExecIn{
			Concurrency:  service.conf.Trigger.Executor.Concurrency,
			ArrangeDelay: service.conf.Trigger.Executor.ArrangeDelay,
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

			in.Triggers[event.Id] = trigger
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(service.conf.Trigger.Executor.Timeout))
		defer cancel()

		out, err := service.uc.Trigger().Exec(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to consume an attempt trigger", "err", err.Error())
			// basically we will not try to retry an attempt trigger
			// because it could be retry later by cronjob
			return map[string]error{}
		}

		if len(out.Error) > 0 {
			// basically we will not try to retry an attempt trigger
			// because it could be retry later by cronjob
			for key, err := range out.Error {
				service.logger.Errorw("consume an attempt trigger got some errors", "key", key, "err", err.Error())
			}
		}

		service.logger.Infow(
			"consumed attempt triggers",
			"event_count", len(events),
			"ok_count", len(out.Success),
			"scheduled", len(out.Scheduled),
			"created", len(out.Created),
			"ko_count", len(out.Error),
		)

		return map[string]error{}
	}
}
