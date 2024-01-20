package cronjob

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/services/recovery/usecase"
)

func UseJob(service *cronjob) func() {
	return func() {
		in := &usecase.ScannerScheduleIn{
			BatchSize: service.conf.Cronjob.BatchSize,
			Buckets:   service.conf.Cronjob.Buckets,
		}
		if err := in.Validate(); err != nil {
			service.logger.Error(err)
			return
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Millisecond*time.Duration(service.conf.Cronjob.Timeout),
		)
		defer cancel()

		out, err := service.uc.Scanner().Schedule(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to execute cronjob", "error", err.Error())
			return
		}
		if len(out.Error) > 0 {
			for ref, err := range out.Error {
				service.logger.Errorw("unable schedule a recovery entities", "ref", ref, "error", err.Error())
			}
		}

		service.logger.Infow("scheduled recovery entities", "ok_count", len(out.Success), "ko_count", len(out.Error))
	}
}
