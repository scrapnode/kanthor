package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc"
)

type TriggerExecIn struct {
	Concurrency int

	ArrangeDelay int64
	Triggers     map[string]*entities.AttemptTrigger
}

func (in *TriggerExecIn) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("concurrency", in.Concurrency, 1),
		validator.NumberGreaterThan("arrange_delay", in.ArrangeDelay, 60000),
		validator.MapRequired("triggers", in.Triggers),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Triggers, func(key string, item *entities.AttemptTrigger) error {
			prefix := fmt.Sprintf("triggers.%s", key)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".app.id", in.Triggers[key].AppId, entities.IdNsApp),
				validator.StringRequired(prefix+".tier", in.Triggers[key].Tier),
				validator.NumberGreaterThan(prefix+".from", in.Triggers[key].From, 0),
				validator.NumberGreaterThan(prefix+".to", int(in.Triggers[key].To), 0),
			)
		}),
	)
}

type TriggerExecOut struct {
	Success []string
	Error   map[string]error

	Scheduled []string
	Created   []string
}

func (uc *trigger) Exec(ctx context.Context, in *TriggerExecIn) (*TriggerExecOut, error) {
	outs := &safe.Slice[*TriggerExecOut]{}

	var wg conc.WaitGroup
	for _, item := range in.Triggers {
		trigger := item
		wg.Go(func() {
			out, err := uc.consume(ctx, trigger, in.Concurrency, in.ArrangeDelay)
			if err != nil {
				uc.logger.Errorw("unable to consume attempt trigger", "trigger", trigger.String(), "err", err.Error())
				return
			}
			outs.Append(out)
		})
	}

	c := make(chan bool)
	defer close(c)

	go func() {
		wg.Wait()
		c <- true
	}()

	select {
	case <-c:
		out := &TriggerExecOut{
			Success:   []string{},
			Error:     map[string]error{},
			Scheduled: []string{},
			Created:   []string{},
		}
		for _, o := range outs.Data() {
			out.merge(o)
		}
		return out, nil
	case <-ctx.Done():
		// actually we may have some success entries, but we can ignore them
		// let cronjob pickup and retry them redundantly
		return nil, ctx.Err()
	}
}

func (uc *trigger) consume(
	ctx context.Context,
	trigger *entities.AttemptTrigger,
	concurrency int,
	delay int64,
) (*TriggerExecOut, error) {
	key := fmt.Sprintf("kanthor.services.attempt.trigger.consumer/%s", trigger.AppId)
	// the lock duration will be long as much as possible
	// so we will time the global timeout as as the lock duration
	// in that duration, we could not consume the same app until the lock is released
	locker := uc.infra.DistributedLockManager(key)

	if err := locker.Lock(ctx); err != nil {
		uc.logger.Errorw("unable to acquire a lock", "key", key)
		// if we could not acquire the lock, don't need to retry so don't set error here
		return nil, err
	}
	defer func() {
		if err := locker.Unlock(ctx); err != nil {
			uc.logger.Errorw("unable to release a lock", "key", key)
		}
	}()

	applicable, err := uc.Applicable(ctx, trigger.AppId)
	if err != nil {
		return nil, err
	}

	from := time.UnixMilli(trigger.From)
	to := time.UnixMilli(trigger.To)

	count, err := uc.repositories.Datastore().Message().Count(ctx, trigger.AppId, from, to)
	if err != nil {
		return nil, err
	}

	out := &TriggerExecOut{
		Success:   []string{},
		Error:     map[string]error{},
		Scheduled: []string{},
		Created:   []string{},
	}

	startdur := time.Now()
	var scanned int64
	ch := uc.repositories.Datastore().Message().Scan(ctx, trigger.AppId, from, to, concurrency)
	for r := range ch {
		scanned += int64(len(r.Data))
		if r.Error != nil {
			return nil, r.Error
		}

		o, err := uc.Perform(ctx, trigger.AppId, r.Data, applicable, delay)
		if err == nil {
			out.merge(o)
		} else {
			uc.logger.Error(err)
		}

		uc.logger.Infof("app_id:%s scanned %.2f%% (%d/%d rows, %v)", trigger.AppId, float64(scanned*100)/float64(count), scanned, count, time.Since(startdur))
	}

	return out, nil
}

func (out *TriggerExecOut) merge(o *TriggerExecOut) {
	out.Success = append(out.Success, o.Success...)
	for k, v := range o.Error {
		out.Error[k] = v
	}
	out.Scheduled = append(out.Scheduled, o.Scheduled...)
	out.Created = append(out.Created, o.Created...)
}
