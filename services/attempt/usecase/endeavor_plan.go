package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndeavorPlanIn struct {
	Timeout int64

	ScanStart int64
	ScanEnd   int64
}

func (in *EndeavorPlanIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("timeout", int(in.Timeout), 1000),
		validator.NumberLessThan("scan_start", in.ScanStart, in.ScanEnd),
	)
}

type EndeavorPlanOut struct {
	Success []string
	From    time.Time
	To      time.Time
}

func (uc *endeavor) Plan(ctx context.Context, in *EndeavorPlanIn) (*EndeavorPlanOut, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(in.Timeout))
	defer cancel()

	from := uc.infra.Timer.Now().Add(time.Duration(in.ScanStart) * time.Millisecond)
	to := uc.infra.Timer.Now().Add(time.Duration(in.ScanEnd) * time.Millisecond)
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
		uc.logger.Debugw("prepare events", "record_count", len(events))

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
		return &EndeavorPlanOut{Success: ok, From: from, To: to}, err
	case <-timeout.Done():
		return &EndeavorPlanOut{Success: ok, From: from, To: to}, timeout.Err()
	}
}

func (uc *endeavor) attempts(ctx context.Context, from, to time.Time) ([]entities.Attempt, error) {
	matching := uc.infra.Timer.Now().UnixMilli()
	attempts, err := uc.repositories.Datastore().Attempt().Scan(ctx, from, to, matching)
	if err != nil {
		return []entities.Attempt{}, err
	}
	uc.logger.Debugw("found records", "record_count", len(attempts))

	returning := []entities.Attempt{}
	for _, attempt := range attempts {
		if attempt.Status == sender.StatusErr {
			returning = append(returning, attempt)
			continue
		}

		if attempt.Status == sender.StatusNone {
			returning = append(returning, attempt)
			continue
		}

		if sender.Is5xxStatus(attempt.Status) {
			returning = append(returning, attempt)
			continue
		}

		uc.logger.Warnw("ignore attempt", "req_id", attempt.ReqId, "status", attempt.Status)
	}

	uc.logger.Debugw("evaluate records", "record_count", len(returning))
	return returning, nil
}
