package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type RetrySelectIn struct {
	BatchSize int
	Triggers  map[string]*entities.AttemptTrigger
}

func ValidateRetrySelectAttemptTrigger(prefix string, attempt *entities.AttemptTrigger) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual(prefix+".to", attempt.From, 0),
		validator.NumberGreaterThan(prefix+".to", attempt.From, attempt.From),
	)
}

func (in *RetrySelectIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("batch_size", in.BatchSize, 0),
		validator.MapRequired("triggers", in.Triggers),
		validator.Map(in.Triggers, func(refId string, item *entities.AttemptTrigger) error {
			prefix := fmt.Sprintf("tasks.%s", refId)
			return ValidateRetrySelectAttemptTrigger(prefix, item)
		}),
	)
}

type RetrySelectOut struct {
	Success []string
	Error   map[string]error
}

func (uc *retry) Select(ctx context.Context, in *RetrySelectIn) (*RetrySelectOut, error) {
	ok := &safe.Map[[]string]{}
	ko := &safe.Map[error]{}

	// we have to store a ref map so if we got any error,
	// we can report back to the caller that a key has a error and should be retry
	eventIdRefs := map[string]string{}
	for eventId, trigger := range in.Triggers {
		eventIdRefs[trigger.String()] = eventId
	}

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		for i := range in.Triggers {
			success, err := uc.choose(ctx, in.Triggers[i], in.BatchSize)
			eventRef := eventIdRefs[in.Triggers[i].String()]

			if err != nil {
				ko.Set(eventRef, err)
				errc <- err
			}

			if len(success) > 0 {
				ok.Set(eventRef, success)
			}
		}

		errc <- nil
	}()

	select {
	case err := <-errc:
		return &RetrySelectOut{Success: ok.Keys(), Error: ko.Data()}, err
	case <-ctx.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, tr := range in.Triggers {
			eventRef := eventIdRefs[tr.String()]

			if _, success := ok.Get(tr.String()); success {
				// already success
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(tr.String()); !has {
				ko.Set(eventRef, ctx.Err())
				continue
			}
		}
		return &RetrySelectOut{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *retry) choose(ctx context.Context, trigger *entities.AttemptTrigger, size int) ([]string, error) {
	query := &entities.ScanningQuery{
		Size: size,
		From: uc.infra.Timer.UnixMilli(trigger.From),
		To:   uc.infra.Timer.UnixMilli(trigger.To),
	}
	ch := uc.repositories.Datastore().Attempt().Scan(ctx, query, uc.infra.Timer.Now().UnixMilli(), uc.conf.Selector.Counter)

	returning := []string{}
	for r := range ch {
		if r.Error != nil {
			return returning, r.Error
		}

		events := map[string]*streaming.Event{}
		for i := range r.Data {
			event, err := transformation.EventFromAttempt(&r.Data[i])
			if err != nil {
				uc.logger.Errorw("ATTEMPT.USECASE.RETRY.SELECT.EVENT_TRANSFORMATION", "error", err.Error(), "attempt", r.Data[i].String())
			}

			events[r.Data[i].ReqId] = event
		}

		errs := uc.publisher.Pub(ctx, events)
		for ref := range events {
			if err, has := errs[ref]; has {
				uc.logger.Errorw("ATTEMPT.USECASE.RETRY.SELECT.EVENT_PUBLISH", "ref", ref, "error", err.Error())
				continue
			}

			returning = append(returning, ref)
		}
	}

	return returning, nil
}
