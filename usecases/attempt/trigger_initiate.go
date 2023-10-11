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

type TriggerInitiateReq struct {
	ScanSize    int
	PublishSize int
}

type TriggerInitiateRes struct {
	Cursor  string
	Success []string
	Error   map[string]error
}

func (uc *trigger) Initiate(ctx context.Context, req *TriggerInitiateReq) (*TriggerInitiateRes, error) {
	apps, cursor, err := uc.applications(ctx, req.ScanSize)
	if err != nil {
		return nil, err
	}
	tiers, err := uc.repos.Application().GetTiers(ctx, apps)
	if err != nil {
		return nil, err
	}

	from := uc.timer.Now().Add(time.Duration(uc.conf.Attempt.Trigger.Consumer.ScanFrom) * time.Second)
	to := uc.timer.Now().Add(time.Duration(uc.conf.Attempt.Trigger.Consumer.ScanTo) * time.Second)

	ok := &ds.SafeSlice[string]{}
	ko := &ds.SafeMap[error]{}

	p := pool.New().WithMaxGoroutines(req.PublishSize)
	for _, app := range apps {
		noti := &entities.AttemptNotification{AppId: app.Id, Tier: tiers[app.Id], From: from.UnixMilli(), To: to.UnixMilli()}

		p.Go(func() {
			event, err := transformation.EventFromNotification(noti)
			if err != nil {
				ko.Set(app.Id, err)
				return
			}

			if err := uc.publisher.Pub(ctx, event); err != nil {
				ko.Set(app.Id, err)

				return
			}

			ko.Set(app.Id, err)
		})
	}
	p.Wait()

	res := &TriggerInitiateRes{Cursor: cursor, Success: ok.Data(), Error: ko.Data()}
	return res, nil
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
