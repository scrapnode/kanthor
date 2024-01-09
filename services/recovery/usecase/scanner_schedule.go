package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services/recovery/config"
)

type ScannerScheduleIn struct {
	BatchSize int
	Buckets   []config.RecoveryScannerBucket
}

func (in *ScannerScheduleIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("batch_size", in.BatchSize, 0),
		validator.Slice(in.Buckets, func(i int, item *config.RecoveryScannerBucket) error {
			return item.Validate(fmt.Sprintf("CONFIG.RECOVERY.SCANNER.BUCKETS[%d]", i))
		}),
	)
}

type ScannerScheduleOut struct {
	Success []string
	Error   map[string]error
}

func (uc *scanner) Schedule(ctx context.Context, in *ScannerScheduleIn) (*ScannerScheduleOut, error) {
	query := &entities.ScanningQuery{
		Limit: in.BatchSize,
	}
	ch := uc.repositories.Database().Application().Scan(ctx, query)

	out := &ScannerScheduleOut{Error: make(map[string]error)}
	for results := range ch {
		if results.Error != nil {
			return nil, results.Error
		}

		events := map[string]*streaming.Event{}
		for _, bucket := range in.Buckets {
			to := uc.infra.Timer.Now().Add(time.Millisecond * time.Duration(-bucket.Offset))
			from := to.Add(time.Millisecond * time.Duration(-bucket.Duration))

			for _, app := range results.Data {
				event, err := transformation.EventFromRecovery(&entities.Recovery{
					Id:    identifier.New(entities.IdNsRec),
					AppId: app.Id,
					To:    to.UnixMilli(),
					From:  from.UnixMilli(),
				})
				if err != nil {
					out.Error[app.Id] = err
					continue
				}

				events[app.Id] = event
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
