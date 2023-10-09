package attempt

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc/pool"
)

type ApplicationTriggerReq struct {
	ScanSize    int
	PublishSize int
}

type ApplicationTriggerRes struct {
	Cursor  string
	Success []string
	Error   map[string]error
}

func (uc *application) Trigger(ctx context.Context, req *ApplicationTriggerReq) (*ApplicationTriggerRes, error) {
	apps, cursor, err := uc.scan(ctx, req.ScanSize)
	if err != nil {
		return nil, err
	}
	tiers, err := uc.repos.Application().GetTiers(ctx, apps)
	if err != nil {
		return nil, err
	}

	ok := &ds.SafeSlice[string]{}
	ko := &ds.SafeMap[error]{}

	p := pool.New().WithMaxGoroutines(req.PublishSize)
	for _, a := range apps {
		app := a
		p.Go(func() {
			err := uc.publish(ctx, &app, tiers[app.Id])
			if err == nil {
				ok.Append(app.Id)
				return
			}
			ko.Set(app.Id, err)
		})
	}
	p.Wait()

	res := &ApplicationTriggerRes{
		Cursor:  cursor,
		Success: ok.Data(),
		Error:   ko.Data(),
	}
	return res, nil
}

func (uc *application) scan(ctx context.Context, size int) ([]entities.Application, string, error) {
	cursor, err := uc.cache.StringGet(ctx, "kanthor.usecases.attempt.application.scan")
	if !errors.Is(err, cache.ErrEntryNotFound) {
		return nil, "", err
	}

	apps, err := uc.repos.Application().Scan(ctx, size, cursor)
	if err != nil {
		return nil, "", err
	}

	if len(apps) > 0 {
		cursor = apps[len(apps)-1].Id
	}

	err = uc.cache.StringSet(ctx, "kanthor.usecases.attempt.application.scan", cursor, time.Hour)
	if err != nil {
		uc.logger.Errorw("unable to set scan cursor to reuse later", "err", err.Error(), "cursor", cursor)
	}

	return apps, "", nil
}

func (uc *application) publish(ctx context.Context, app *entities.Application, tier string) error {
	event, err := transformation.EventFromApplication(app, tier)
	if err != nil {
		return err
	}

	if err := uc.publisher.Pub(ctx, event); err != nil {
		return err
	}

	return nil
}
