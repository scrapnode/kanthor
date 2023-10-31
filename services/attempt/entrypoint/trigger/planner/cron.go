package planner

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

// @TODO: remove hardcode
var key = "attempt.trigger.cron"

func RegisterCron(service *planner) func() {
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
		if err := ucreq.Validate(); err != nil {
			service.logger.Errorw("invalid trigger plan request", "err", err.Error(), "ucreq", utils.Stringify(ucreq))
			return
		}

		ucres, err := service.uc.Trigger().Plan(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to plan attempt triggers", "err", err.Error())
			return
		}

		service.logger.Infow("planned attempt triggers", "count", len(ucres.Success), "from", ucres.From.Format(time.RFC3339), "to", ucres.To.Format(time.RFC3339))
	}
}
