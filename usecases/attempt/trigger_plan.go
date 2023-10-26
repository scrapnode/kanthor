package attempt

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

type TriggerPlanReq struct {
	Size    int
	Timeout int64

	ScanStart int64
	ScanEnd   int64
}

func (req *TriggerPlanReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("size", req.Size, 1),
		validator.NumberGreaterThan("timeout", int(req.Timeout), 1000),
		validator.NumberGreaterThan("scan_start", req.ScanStart, req.ScanEnd),
		validator.NumberLessThan("scan_end", req.ScanEnd, req.ScanStart),
	)
}

type TriggerPlanRes struct {
	Success []string
}

func (uc *trigger) Plan(ctx context.Context, req *TriggerPlanReq) (*TriggerPlanRes, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(req.Timeout))
	defer cancel()

	ok := []string{}

	errc := make(chan error)
	defer close(errc)

	go func() {
		apps, err := uc.applications(ctx, req.Size)
		if err != nil {
			errc <- err
			return
		}
		tiers, err := uc.repos.Application().GetTiers(ctx, apps)
		if err != nil {
			errc <- err
			return
		}

		from := uc.infra.Timer.Now().Add(time.Duration(req.ScanStart) * time.Millisecond)
		to := uc.infra.Timer.Now().Add(time.Duration(req.ScanEnd) * time.Millisecond)

		events := map[string]*streaming.Event{}
		for _, app := range apps {
			key := app.Id
			trigger := &entities.AttemptTrigger{
				AppId: app.Id,
				Tier:  tiers[app.Id],
				From:  from.UnixMilli(),
				To:    to.UnixMilli(),
			}
			event, err := transformation.EventFromTrigger(trigger)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform trigger to event", "trigger", trigger.String())
				continue
			}
			events[key] = event
		}

		var perr error
		errs := uc.infra.Stream.Publisher("attempt_trigger_plan").Pub(ctx, events)
		for key := range events {
			if err, ok := errs[key]; ok {
				perr = errors.Join(perr, err)
			}

			ok = append(ok, key)
		}

		errc <- nil
	}()

	select {
	case err := <-errc:
		return &TriggerPlanRes{Success: ok}, err
	case <-timeout.Done():
		return &TriggerPlanRes{Success: ok}, timeout.Err()
	}
}

// @TODO: remove it
var key = "kanthor.usecases.attempt.trigger.scan"

func (uc *trigger) applications(ctx context.Context, size int) ([]entities.Application, error) {
	cursor, err := uc.infra.Cache.StringGet(ctx, key)
	// only accept entry not found error
	if err != nil && !errors.Is(err, cache.ErrEntryNotFound) {
		return nil, err
	}

	apps, err := uc.repos.Application().Scan(ctx, size, cursor)
	if err != nil {
		return nil, err
	}

	nomore := len(apps) == 0 && cursor != ""
	if nomore {
		uc.logger.Warnw("no more applications to scan", "cursor", cursor)
		// no more app, reset cursor
		uc.infra.Cache.Del(ctx, key)
		return []entities.Application{}, nil
	}

	if len(apps) > 0 {
		cursor = apps[len(apps)-1].Id
	}

	err = uc.infra.Cache.StringSet(ctx, key, cursor, time.Hour)
	if err != nil {
		uc.logger.Errorw("unable to set scan cursor to reuse later", "err", err.Error(), "cursor", cursor)
	}

	return apps, nil
}
