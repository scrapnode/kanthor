package planner

import (
	"context"

	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func RegisterCron(service *planner) func() {
	key := "kanthor.services.attempt.trigger"

	return func() {
		locker := service.infra.DistributedLockManager(key)
		ctx, cancel := context.WithTimeout(context.Background(), locker.TimeToLive())
		defer cancel()

		if err := locker.Lock(ctx); err != nil {
			service.logger.Errorw("unable to acquire a lock", "key", key, "err", err.Error())
			return
		}
		defer func() {
			if err := locker.Unlock(ctx); err != nil {
				service.logger.Errorw("unable to release a lock", "key", key, "err", err.Error())
			}
		}()

		ucreq := &usecase.TriggerPlanReq{
			Timeout:   service.conf.Trigger.Planner.Timeout,
			Size:      service.conf.Trigger.Planner.Size,
			ScanStart: service.conf.Trigger.Planner.ScanStart,
			ScanEnd:   service.conf.Trigger.Planner.ScanEnd,
		}
		ucres, err := service.uc.Trigger().Plan(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to plan attempt trigger", "err", err.Error())
			return
		}

		service.logger.Infow("planned attempt triggers", "count", len(ucres.Success))
	}
}
