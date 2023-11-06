package planner

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

// @TODO: remove hardcode
var key = "attempt.endeavor.cron"

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

		in := &usecase.EndeavorPlanIn{
			Timeout:   service.conf.Endeavor.Planner.Timeout,
			ScanStart: service.conf.Endeavor.Planner.ScanStart,
			ScanEnd:   service.conf.Endeavor.Planner.ScanEnd,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw("invalid trigger plan request", "err", err.Error(), "in", utils.Stringify(in))
			return
		}

		out, err := service.uc.Endeavor().Plan(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to plan attempt endeavors", "err", err.Error())
			return
		}

		service.logger.Infow("planned attempt endeavors", "count", len(out.Success), "from", out.From.Format(time.RFC3339), "to", out.To.Format(time.RFC3339))
	}
}
