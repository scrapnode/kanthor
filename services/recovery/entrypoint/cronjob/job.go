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
			return
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Millisecond*time.Duration(service.conf.Cronjob.Timeout),
		)
		defer cancel()

		_, _ = service.uc.Scanner().Schedule(ctx, in)
	}
}
