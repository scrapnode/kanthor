package attempt

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc/pool"
)

type TriggerPlanReq struct {
	ScanFrom int64
	ScanTo   int64

	ChunkTimeout int64
	ChunkSize    int
}

func (req *TriggerPlanReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("scan_from", req.ScanFrom, req.ScanTo),
		validator.NumberLessThan("scan_to", req.ScanTo, req.ScanFrom),
		validator.NumberGreaterThan("chunk_timeout", req.ChunkTimeout, 1000),
		validator.NumberGreaterThan("chunk_size", req.ChunkSize, 1),
	)
}

type TriggerPlanRes struct {
	Cursor  string
	Success []string
	Error   map[string]error
}

func (uc *trigger) Plan(ctx context.Context, req *TriggerPlanReq) (*TriggerPlanRes, error) {
	apps, cursor, err := uc.applications(ctx, req.ChunkSize)
	if err != nil {
		return nil, err
	}
	tiers, err := uc.repos.Application().GetTiers(ctx, apps)
	if err != nil {
		return nil, err
	}

	from := uc.timer.Now().Add(time.Duration(req.ScanFrom) * time.Millisecond)
	to := uc.timer.Now().Add(time.Duration(req.ScanTo) * time.Millisecond)

	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}

	// timeout duration will be scaled based on how many applications you have
	duration := time.Duration(req.ChunkTimeout * int64(len(apps)+1))
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*duration)
	defer cancel()

	p := pool.New().WithMaxGoroutines(req.ChunkSize)
	for _, app := range apps {
		notification := &entities.AttemptNotification{AppId: app.Id, Tier: tiers[app.Id], From: from.UnixMilli(), To: to.UnixMilli()}

		p.Go(func() {
			event, err := transformation.EventFromNotification(notification)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform notification to event", "notification", notification.String())
				return
			}

			if err := uc.publisher.Pub(ctx, event); err != nil {
				ko.Set(app.Id, err)

				return
			}

			ko.Set(app.Id, err)
		})
	}

	c := make(chan bool)
	defer close(c)

	go func() {
		p.Wait()
		c <- true
	}()

	select {
	case <-c:
		return &TriggerPlanRes{Cursor: cursor, Success: ok.Data(), Error: ko.Data()}, nil
	case <-timeout.Done():
		// we don't need to check which notication was consumed
		// because it could be simply retry later with the cronjob
		// once the context deadline is exceeded, consider all notications are failed
		return nil, ctx.Err()
	}
}

func (uc *trigger) applications(ctx context.Context, size int) ([]entities.Application, string, error) {
	cursor, err := uc.cache.StringGet(ctx, "kanthor.usecases.attempt.trigger.scan")
	if !errors.Is(err, cache.ErrEntryNotFound) {
		return nil, "", err
	}

	apps, err := uc.repos.Application().Scan(ctx, size, cursor)
	if err != nil {
		return nil, "", err
	}

	// if we scanned with a cursor and there is no app,
	// we should retry without the cursor so we can start scanning at the beginning of the dataset
	backToBeginning := len(apps) == 0 && cursor != ""
	if backToBeginning {
		apps, err = uc.repos.Application().Scan(ctx, size, "")
		if err != nil {
			return nil, "", err
		}
	}

	if len(apps) > 0 {
		cursor = apps[len(apps)-1].Id
	}

	err = uc.cache.StringSet(ctx, "kanthor.usecases.attempt.trigger.scan", cursor, time.Hour)
	if err != nil {
		uc.logger.Errorw("unable to set scan cursor to reuse later", "err", err.Error(), "cursor", cursor)
	}

	return apps, "", nil
}
