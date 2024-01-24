package trigger

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func UseJob(service *trigger) func() {
	return func() {
		in := &usecase.RetryTriggerIn{
			Buckets: service.conf.Cronjob.Buckets,
		}
		if err := in.Validate(); err != nil {
			return
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Millisecond*time.Duration(service.conf.Cronjob.Timeout),
		)
		defer cancel()

		_, _ = service.uc.Retry().Trigger(ctx, in)
	}
}
