package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ScannerExecuteIn struct {
	RecoveryBatchSize int
	Tasks             map[string]*entities.RecoveryTask
}

func ValidateRecoveryTask(prefix string, recovery *entities.RecoveryTask) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".app_id", recovery.AppId, entities.IdNsApp),
		validator.NumberGreaterThanOrEqual(prefix+".to", recovery.From, 0),
		validator.NumberGreaterThan(prefix+".to", recovery.From, recovery.From),
	)
}

func (in *ScannerExecuteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("recovery_batch_size", in.RecoveryBatchSize, 0),
		validator.MapRequired("tasks", in.Tasks),
		validator.Map(in.Tasks, func(refId string, item *entities.RecoveryTask) error {
			prefix := fmt.Sprintf("tasks.%s", refId)
			return ValidateRecoveryTask(prefix, item)
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

	// we have to store a ref map of recovery.app_id and the key
	// so if we got any error, we can report back to the call that a key has a error
	eventIdRefs := map[string]string{}
	for eventId, task := range in.Tasks {
		eventIdRefs[task.AppId] = eventId
	}

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		for i := range in.Tasks {
			success, err := uc.execute(ctx, in.Tasks[i], in.RecoveryBatchSize)
			eventRef := eventIdRefs[in.Tasks[i].AppId]

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
		for _, rec := range in.Tasks {
			eventRef := eventIdRefs[rec.AppId]

			if _, success := ok.Get(rec.AppId); success {
				// already success
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(rec.AppId); !has {
				ko.Set(eventRef, ctx.Err())
				continue
			}
		}
		return &ScannerExecuteOut{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *scanner) execute(ctx context.Context, recovery *entities.RecoveryTask, size int) ([]string, error) {
	routes, err := uc.repositories.Database().Application().GetRoutes(ctx, []string{recovery.AppId})
	if err != nil {
		return nil, err
	}
	if _, has := routes[recovery.AppId]; !has {
		return []string{}, nil
	}

	query := &entities.ScanningQuery{
		Size: size,
		From: uc.infra.Timer.UnixMilli(recovery.From),
		To:   uc.infra.Timer.UnixMilli(recovery.To),
	}
	ch := uc.repositories.Datastore().Message().Scan(ctx, recovery.AppId, query)

	returning := []string{}
	for r := range ch {
		if r.Error != nil {
			return returning, r.Error
		}

		conditions := make([]string, 0)
		messages := map[string]*entities.Message{}
		routeMaps := map[string]*routing.Route{}

		for i := range r.Data {
			route, has := routes[r.Data[i].AppId]
			if !has {
				continue
			}

			for j := range route {
				pair := fmt.Sprintf("%s/%s", route[j].Endpoint.Id, r.Data[i].Id)
				conditions = append(conditions, pair)
				messages[pair] = &r.Data[i]
				routeMaps[pair] = &route[j]
			}
		}

		scheduled, err := uc.repositories.Datastore().Request().Check(ctx, conditions)
		if err != nil {
			return returning, err
		}

		events := map[string]*streaming.Event{}
		for pair, ok := range scheduled {
			// already scheduled a request for the message to the endpoint, ignore
			if ok {
				continue
			}

			request, trace := routing.PlanRequest(uc.infra.Timer, messages[pair], routeMaps[pair])

			if request != nil {
				if event, err := transformation.EventFromRequest(request); err == nil {
					events[pair] = event
				} else {
					uc.logger.Errorw("RECOVERY.USECASE.SCANNER.EXECUTE.EVENT_TRANSFORMATION", "error", err.Error())
				}
			}

			if len(trace) > 0 {
				uc.logger.Warnw(trace[0].(string), trace[1:]...)
			}
		}

		errs := uc.publisher.Pub(ctx, events)
		for ref := range events {
			if err, has := errs[ref]; has {
				uc.logger.Errorw("RECOVERY.USECASE.SCANNER.EXECUTE.EVENT_PUBLISH", "ref", ref, "error", err.Error())
				continue
			}

			returning = append(returning, ref)
		}
	}

	return returning, nil
}
