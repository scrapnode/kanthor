package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndeavorPlanReq struct {
	Timeout int64

	ScanStart int64
	ScanEnd   int64
}

func (req *EndeavorPlanReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("timeout", int(req.Timeout), 1000),
		validator.NumberGreaterThan("scan_start", req.ScanStart, req.ScanEnd),
		validator.NumberLessThan("scan_end", req.ScanEnd, req.ScanStart),
	)
}

type EndeavorPlanRes struct {
	Success []string
	From    time.Time
	To      time.Time
}

func (uc *endeavor) Plan(ctx context.Context, req *EndeavorPlanReq) (*EndeavorPlanRes, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(req.Timeout))
	defer cancel()

	from := uc.infra.Timer.Now().Add(time.Duration(req.ScanStart) * time.Millisecond)
	to := uc.infra.Timer.Now().Add(time.Duration(req.ScanEnd) * time.Millisecond)
	ok := []string{}

	errc := make(chan error)
	defer close(errc)
	go func() {
		atts, err := uc.attempts(ctx, from, to)
		if err != nil {
			errc <- err
			return
		}

		events := map[string]*streaming.Event{}
		for _, att := range atts {
			event, err := transformation.EventFromAttempt(&att)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform attempt to event", "attempt", att.String())
				continue
			}
			events[key] = event
		}

		var perr error
		errs := uc.infra.Stream.Publisher("attempt_endeavor_plan").Pub(ctx, events)
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
		return &EndeavorPlanRes{Success: ok, From: from, To: to}, err
	case <-timeout.Done():
		return &EndeavorPlanRes{Success: ok, From: from, To: to}, timeout.Err()
	}
}

func (uc *endeavor) attempts(ctx context.Context, from, to time.Time) ([]entities.Attempt, error) {
	matching := uc.infra.Timer.Now().UnixMilli()
	attempts, err := uc.repositories.Datastore().Attempt().Scan(ctx, from, to, matching)
	if err != nil {
		return []entities.Attempt{}, err
	}

	returning := []entities.Attempt{}
	for _, attempt := range attempts {
		// @TODO: remove hardcode
		if attempt.Status == 0 {
			attempts = append(attempts, attempt)
			continue
		}

		if entities.Is5xx(attempt.Status) {
			attempts = append(attempts, attempt)
			continue
		}

		uc.logger.Warnw("ignore attempt", "req_id", attempt.ReqId, "status", attempt.Status)
	}
	return returning, nil
}
