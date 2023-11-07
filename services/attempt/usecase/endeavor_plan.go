package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/status"
	"github.com/scrapnode/kanthor/domain/transformation"
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
		s, err := uc.scan(ctx, from, to)
		if err != nil {
			errc <- err
			return
		}

		ids := uc.trigger(ctx, s)
		ok = append(ok, ids...)

		if err := uc.repositories.Datastore().Attempt().MarkIgnore(ctx, s.Ignore); err != nil {
			uc.logger.Errorw("unable to ignore attempts", "req_ids", s.Ignore)
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

type strive struct {
	Attemptable []entities.Attempt
	Ignore      []string
}

func (uc *endeavor) scan(ctx context.Context, from, to time.Time) (*strive, error) {
	returning := &strive{Attemptable: []entities.Attempt{}, Ignore: []string{}}

	less := uc.infra.Timer.Now().UnixMilli()
	attempts, err := uc.repositories.Datastore().Attempt().Scan(ctx, from, to, less)
	if err != nil {
		return nil, err
	}
	uc.logger.Debugw("found records", "record_count", len(attempts))

	for _, attempt := range attempts {
		// ignore
		if attempt.Status == status.ErrIgnore {
			continue
		}
		// or Ignore
		if attempt.Complete() {
			continue
		}

		if attempt.Status == status.ErrUnknown {
			returning.Attemptable = append(returning.Attemptable, attempt)
			continue
		}

		if attempt.Status == status.None {
			returning.Attemptable = append(returning.Attemptable, attempt)
			continue
		}

		if status.Is5xx(attempt.Status) {
			returning.Attemptable = append(returning.Attemptable, attempt)
			continue
		}

		returning.Ignore = append(returning.Ignore, attempt.ReqId)
		uc.logger.Warnw("ignore attempt", "req_id", attempt.ReqId, "status", attempt.Status)
	}

	uc.logger.Debugw("evaluate records", "attempt_count", len(returning.Attemptable), "Ignore_count", len(returning.Ignore))
	return returning, nil
}

func (uc *endeavor) trigger(ctx context.Context, s *strive) []string {
	events := map[string]*streaming.Event{}
	for _, att := range s.Attemptable {
		refId := att.ReqId
		event, err := transformation.EventFromAttempt(&att)
		if err != nil {
			// un-recoverable error
			uc.logger.Errorw("could not transform attempt to event", "attempt", att.String())
			continue
		}
		events[refId] = event
	}

	errs := uc.infra.Stream.Publisher("attempt_endeavor_plan").Pub(ctx, events)
	ok := []string{}

	for refId := range events {
		if err, ok := errs[refId]; ok {
			uc.logger.Errorw("trigger event got error", "req_id", refId, "error", err.Error())
			continue
		}

		ok = append(ok, refId)
	}

	return ok
}
