package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services/attempt/config"
)

type RetryTriggerIn struct {
	Buckets []config.AttemptBucket
}

func (in *RetryTriggerIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.Slice(in.Buckets, func(i int, item *config.AttemptBucket) error {
			return item.Validate(fmt.Sprintf("buckets[%d]", i))
		}),
	)
}

type RetryTriggerOut struct {
	Success []string
	Error   map[string]error
}

func (uc *retry) Trigger(ctx context.Context, in *RetryTriggerIn) (*RetryTriggerOut, error) {
	out := &RetryTriggerOut{Error: make(map[string]error)}

	now := uc.infra.Timer.Now()
	events := map[string]*streaming.Event{}
	for _, bucket := range in.Buckets {
		key := fmt.Sprintf("%d/%d", bucket.Offset, bucket.Duration)
		to := now.Add(time.Millisecond * time.Duration(-bucket.Offset))
		from := to.Add(time.Millisecond * time.Duration(-bucket.Duration))

		event, err := transformation.EventFromAttemptTrigger(&entities.AttemptTrigger{
			To:   to.UnixMilli(),
			From: from.UnixMilli(),
			// this is a hidden feature and is convenient way to punish the tasks and bypass deduplicated logic
			// Init:  uc.infra.Timer.Now().UnixMilli(),
		})
		if err != nil {
			out.Error[key] = err
			continue
		}

		events[key] = event
	}

	errs := uc.publisher.Pub(ctx, events)
	for key, err := range errs {
		out.Error[key] = err
	}

	for refId := range events {
		if _, exist := out.Error[refId]; !exist {
			out.Success = append(out.Success, refId)
		}
	}

	return out, nil
}
