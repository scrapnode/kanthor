package attempt

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/attempt"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

func RegisterTriggerCron(service *attempt) func() {
	key := "kanthor.services.attempt.trigger"
	duration := time.Duration(service.conf.Attempt.Trigger.Cron.LockDuration) * time.Second

	return func() {
		locker := service.locker(key, duration)
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		if err := locker.Lock(ctx); err != nil {
			service.logger.Errorw("unable to acquire a lock", "key", key)
			return
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			if err := locker.Unlock(ctx); err != nil {
				service.logger.Errorw("unable to release a lock", "key", key)
			}
		}()

		ucreq := &usecase.TriggerInitiateReq{
			ScanSize:    service.conf.Attempt.Trigger.Cron.ScanSize,
			PublishSize: service.conf.Attempt.Trigger.Cron.PublishSize,
		}
		ucres, err := service.uc.Trigger().Initiate(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to initiate attempt notification", "err", err.Error())
			return
		}

		if len(ucres.Error) > 0 {
			for key, err := range ucres.Error {
				service.logger.Errorw("initiate attempt notification got err", "key", key, "err", err.Error())
			}
		}
	}
}

func RegisterTriggerConsumer(service *attempt) streaming.SubHandler {
	return func(events []*streaming.Event) map[string]error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		notifications := []entities.AttemptNotification{}
		for _, event := range events {
			notification, err := transformation.EventToNotification(event)
			if err != nil {
				service.logger.Errorw("unable to transform event to attempt notification", "err", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			notifications = append(notifications, *notification)
		}

		ucreq := &usecase.TriggerConsumeReq{ConsumeSize: 1, Notifications: notifications}
		ucres, err := service.uc.Trigger().Consume(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to consume attempt notifications", "err", err.Error())
			// basically we will not try to retry an attempt notification
			// because it could be retry later by cronjob
			return map[string]error{}
		}

		if len(ucres.Error) > 0 {
			for key, err := range ucres.Error {
				service.logger.Errorw("initiate attempt notification got err", "key", key, "err", err.Error())
			}
		}

		return map[string]error{}
	}
}
