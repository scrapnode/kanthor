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
	duration := time.Millisecond * time.Duration(service.conf.Attempt.Trigger.Plan.LockDuration)

	return func() {
		locker := service.infra.DistributedLockManager(key, duration)
		ctx := context.Background()

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

		ucreq := &usecase.TriggerPlanReq{
			Timeout:   service.conf.Attempt.Trigger.Plan.Timeout,
			RateLimit: service.conf.Attempt.Trigger.Plan.RateLimit,
			ScanStart: service.conf.Attempt.Trigger.Plan.ScanStart,
			ScanEnd:   service.conf.Attempt.Trigger.Plan.ScanEnd,
		}
		ucres, err := service.uc.Trigger().Plan(ctx, ucreq)
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
		ucreq := &usecase.TriggerExecReq{
			Timeout:       service.conf.Attempt.Trigger.Exec.Timeout,
			RateLimit:     service.conf.Attempt.Trigger.Exec.RateLimit,
			Delay:         service.conf.Attempt.Trigger.Exec.Delay,
			Notifications: []entities.AttemptNotification{},
		}

		for _, event := range events {
			notification, err := transformation.EventToNotification(event)
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
