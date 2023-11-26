package planner

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/scrapnode/kanthor/infrastructure/dlm"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

// @TODO: remove hardcode
var key = "attempt.endeavor.cron"

func Handler(service *planner, schedule *cron.SpecSchedule) func() {
	return func() {
		ttl := schedule.Second + schedule.Minute + schedule.Hour + schedule.Dom + schedule.Month + schedule.Dow
		// lock longger than timeout
		locker := service.infra.DistributedLockManager(key, dlm.TimeToLive(ttl+uint64(time.Minute)))

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ttl))
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
			Size:      service.conf.Endeavor.Planner.Size,
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
		service.logger.Infow("waiting for next schedule", "next_scheule", schedule.Next(time.Now().UTC()).Format(time.RFC3339))
	}
}
