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

type ScannerScheduleIn struct {
	BatchSize int
	Buckets   []config.AttemptBucket
}

func (in *ScannerScheduleIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("batch_size", in.BatchSize, 0),
		validator.Slice(in.Buckets, func(i int, item *config.AttemptBucket) error {
			return item.Validate(fmt.Sprintf("buckets[%d]", i))
		}),
	)
}

type ScannerScheduleOut struct {
	Success []string
	Error   map[string]error
}

func (uc *scanner) Schedule(ctx context.Context, in *ScannerScheduleIn) (*ScannerScheduleOut, error) {
	query := &entities.ScanningQuery{
		Size: in.BatchSize,
	}
	ch := uc.repositories.Database().Endpoint().Scan(ctx, query)

	out := &ScannerScheduleOut{Error: make(map[string]error)}
	for results := range ch {
		if results.Error != nil {
			return nil, results.Error
		}

		now := uc.infra.Timer.Now()
		events := map[string]*streaming.Event{}
		for _, bucket := range in.Buckets {
			to := now.Add(time.Millisecond * time.Duration(-bucket.Offset))
			from := to.Add(time.Millisecond * time.Duration(-bucket.Duration))

			for _, ep := range results.Data {
				event, err := transformation.EventFromAttemptTask(&entities.AttemptTask{
					AppId: ep.AppId,
					EpId:  ep.Id,
					To:    to.UnixMilli(),
					From:  from.UnixMilli(),
					// this is a hidden feature and is convenient way to punish the tasks and bypass deduplicated logic
					// Init:  uc.infra.Timer.Now().UnixMilli(),
				})
				if err != nil {
					out.Error[ep.Id] = err
					continue
				}

				events[ep.Id] = event
			}
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
	}

	return out, nil
}
