package trigger

import (
	"context"

	usecase "github.com/scrapnode/kanthor/usecases/attempt"
)

func RegisterCron(service *planner) func() {
	key := "kanthor.services.attempt.trigger"

	return func() {
		locker := service.infra.DistributedLockManager(key)
		ctx, cancel := context.WithDeadline(context.Background(), locker.Until())
		defer cancel()

		if err := locker.Lock(ctx); err != nil {
			service.logger.Errorw("unable to acquire a lock", "key", key)
			return
		}
		defer func() {
			if err := locker.Unlock(ctx); err != nil {
				service.logger.Errorw("unable to release a lock", "key", key)
			}
		}()

		ucreq := &usecase.TriggerPlanReq{
			Timeout:   service.conf.Attempt.Trigger.Planner.Timeout,
			RateLimit: service.conf.Attempt.Trigger.Planner.RateLimit,
			ScanStart: service.conf.Attempt.Trigger.Planner.ScanStart,
			ScanEnd:   service.conf.Attempt.Trigger.Planner.ScanEnd,
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
