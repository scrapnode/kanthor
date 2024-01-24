package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/status"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ScannerExecuteIn struct {
	BatchSize int
	Tasks     map[string]*entities.AttemptTask
}

func ValidateScannerExecuteAttemptTask(prefix string, attempt *entities.AttemptTask) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".ep_id", attempt.EpId, entities.IdNsEp),
		validator.NumberGreaterThanOrEqual(prefix+".to", attempt.From, 0),
		validator.NumberGreaterThan(prefix+".to", attempt.From, attempt.From),
	)
}

func (in *ScannerExecuteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("batch_size", in.BatchSize, 0),
		validator.MapRequired("tasks", in.Tasks),
		validator.Map(in.Tasks, func(refId string, item *entities.AttemptTask) error {
			prefix := fmt.Sprintf("tasks.%s", refId)
			return ValidateScannerExecuteAttemptTask(prefix, item)
		}),
	)
}

type ScannerExecuteOut struct {
	Success []string
	Error   map[string]error
}

func (uc *scanner) Execute(ctx context.Context, in *ScannerExecuteIn) (*ScannerExecuteOut, error) {
	ok := &safe.Map[[]string]{}
	ko := &safe.Map[error]{}

	// we have to store a ref map so if we got any error,
	// we can report back to the caller that a key has a error and should be retry
	eventIdRefs := map[string]string{}
	for eventId, task := range in.Tasks {
		eventIdRefs[task.EpId] = eventId
	}

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		for i := range in.Tasks {
			success, err := uc.execute(ctx, in.Tasks[i], in.BatchSize)
			eventRef := eventIdRefs[in.Tasks[i].EpId]

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
		return &ScannerExecuteOut{Success: ok.Keys(), Error: ko.Data()}, err
	case <-ctx.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, task := range in.Tasks {
			eventRef := eventIdRefs[task.EpId]

			if _, success := ok.Get(task.EpId); success {
				// already success
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(task.EpId); !has {
				ko.Set(eventRef, ctx.Err())
				continue
			}
		}
		return &ScannerExecuteOut{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *scanner) execute(ctx context.Context, attempt *entities.AttemptTask, size int) ([]string, error) {
	query := &entities.ScanningQuery{
		Size: size,
		From: uc.infra.Timer.UnixMilli(attempt.From),
		To:   uc.infra.Timer.UnixMilli(attempt.To),
	}
	ch := uc.repositories.Datastore().Request().Scan(ctx, attempt.EpId, query)

	returning := []string{}
	for r := range ch {
		if r.Error != nil {
			return returning, r.Error
		}

		requests := map[string]*entities.Request{}
		var msgIds []string
		for i := range r.Data {
			msgIds = append(msgIds, r.Data[i].MsgId)
			// because we are working in the same endpoint so each message will be unique
			// that why we can use msg_id as a key of request reference
			requests[r.Data[i].MsgId] = &r.Data[i]
		}

		responses, err := uc.repositories.Datastore().Response().Check(ctx, attempt.EpId, msgIds)
		if err != nil {
			return returning, err
		}

		now := uc.infra.Timer.Now()
		events := map[string]*streaming.Event{}
		for msgId, state := range responses {
			// already scheduled a request for the message to the endpoint, ignore
			if status.IsAnyOK(state) {
				continue
			}

			attempt := &entities.Attempt{
				ReqId: requests[msgId].Id,
				MsgId: msgId,
				EpId:  requests[msgId].EpId,
				AppId: requests[msgId].AppId,
				Tier:  requests[msgId].Tier,
			}
			attempt.ScheduleNext = now.Add(time.Millisecond * time.Duration(uc.conf.Consumer.ScheduleDelay)).UnixMilli()
			attempt.ScheduledAt = now.UnixMilli()

			event, err := transformation.EventFromAttempt(attempt)
			if err != nil {
				uc.logger.Errorw("ATTEMPT.USECASE.SCANNER.EXECUTE.EVENT_TRANSFORMATION", "error", err.Error(), "attempt", attempt.String())
			}

			events[msgId] = event
		}

		errs := uc.publisher.Pub(ctx, events)
		for ref := range events {
			if err, has := errs[ref]; has {
				uc.logger.Errorw("ATTEMPT.USECASE.SCANNER.EXECUTE.EVENT_PUBLISH", "ref", ref, "error", err.Error())
				continue
			}

			returning = append(returning, ref)
		}
	}

	return returning, nil
}
