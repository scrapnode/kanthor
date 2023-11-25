package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/assessor"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/status"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services/attempt/repositories/ds"
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
			out.Success = append(out.Success, o.Success...)
			for k, v := range o.Error {
				out.Error[k] = v
			}
			out.Scheduled = append(out.Scheduled, o.Scheduled...)
			out.Created = append(out.Created, o.Created...)
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

	applicable, err := uc.applicable(ctx, trigger.AppId)
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

	var scanned int64
	ch := uc.repositories.Datastore().Message().Scan(ctx, trigger.AppId, from, to, concurrency)
	for r := range ch {
		scanned += int64(len(r.Data))
		uc.logger.Infof("app_id:%s scanned %d/%d rows", trigger.AppId, scanned, count)
		if r.Error != nil {
			return nil, r.Error
		}

		msgIds, requests, err := uc.examine(ctx, trigger.AppId, applicable, r.Data)
		if err != nil {
			return nil, err
		}

		scheduledok, scheduledko := uc.schedule(ctx, applicable, msgIds)
		out.Success = append(out.Success, scheduledok.Data()...)
		for k, v := range scheduledko.Data() {
			out.Error[k] = v
		}
		out.Scheduled = append(out.Scheduled, scheduledok.Data()...)

		createdok, createdko := uc.create(ctx, delay, requests)
		out.Success = append(out.Success, createdok.Data()...)
		for k, v := range createdko.Data() {
			out.Error[k] = v
		}
		out.Scheduled = append(out.Scheduled, createdok.Data()...)
	}

	return out, nil
}

// examine messages, requests and responses
func (uc *trigger) examine(
	ctx context.Context,
	appId string,
	applicable *assessor.Assets,
	messages map[string]entities.Message,
) ([]entities.Message, []entities.Request, error) {
	var msgIds []string
	for id := range messages {
		msgIds = append(msgIds, id)
	}

	requests, err := uc.repositories.Datastore().Request().Scan(ctx, appId, msgIds)
	if err != nil {
		return nil, nil, err
	}

	// got duplicate rows
	responses, err := uc.repositories.Datastore().Response().Scan(ctx, appId, msgIds)
	if err != nil {
		return nil, nil, err
	}

	schedulable := []entities.Message{}
	attemptable := []entities.Request{}

	status := uc.hash(requests, responses)
	for _, message := range messages {
		for _, ep := range applicable.EndpointMap {
			reqId, hasReq := status[ds.ReqKey(message.Id, ep.Id)]
			if !hasReq {
				// no request -> must schedule message again -> don't create any attempt
				schedulable = append(schedulable, message)
				continue
			}

			_, hasRes := status[ds.ResKey(message.Id, ep.Id)]
			if !hasRes {
				// has request + no success response -> create an attempt
				attemptable = append(attemptable, requests[reqId])
				continue
			}

			// has success response, ignore
		}
	}

	return schedulable, attemptable, nil
}

func (uc *trigger) applicable(ctx context.Context, appId string) (*assessor.Assets, error) {
	key := utils.Key("scheduler", appId)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour, func() (*assessor.Assets, error) {
		endpoints, err := uc.repositories.Database().Endpoint().List(ctx, appId)
		if err != nil {
			return nil, err
		}
		returning := &assessor.Assets{EndpointMap: map[string]entities.Endpoint{}}
		for _, ep := range endpoints {
			returning.EndpointMap[ep.Id] = ep
		}

		rules, err := uc.repositories.Database().Endpoint().Rules(ctx, appId)
		if err != nil {
			return nil, err
		}
		returning.Rules = rules

		return returning, nil
	})
}

func (uc *trigger) hash(requests map[string]entities.Request, responses map[string]ds.ResponseStatusRow) map[string]string {
	returning := map[string]string{}

	for _, request := range requests {
		// for checking whether we have scheduled a request for an endpoint or not
		// if no request was scheduled, we should schedule it instead of create an attempt
		key := ds.ReqKey(request.MsgId, request.EpId)
		returning[key] = request.Id
	}

	for _, response := range responses {
		key := ds.ResKey(response.MsgId, response.EpId)

		// we already recognized that the endpoint had success status, don't need to check any more
		if _, has := returning[key]; has {
			continue
		}

		// status is ok, saved the success response id
		if status.Is2xx(response.Status) {
			returning[key] = response.Id
		}
	}

	return returning
}

// schedule messages
func (uc *trigger) schedule(ctx context.Context, applicable *assessor.Assets, messages []entities.Message) (*safe.Slice[string], *safe.Map[error]) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}
	requests := []entities.Request{}

	for _, message := range messages {
		ins, logs := assessor.Requests(&message, applicable, uc.infra.Timer)
		if len(logs) > 0 {
			for _, l := range logs {
				uc.logger.Warnw(l[0].(string), l[1:]...)
			}
		}
		requests = append(requests, ins...)
	}

	if len(requests) == 0 {
		return ok, ko
	}

	events := map[string]*streaming.Event{}
	for _, request := range requests {
		key := utils.Key(request.AppId, request.MsgId, request.EpId, request.Id)
		event, err := transformation.EventFromRequest(&request)
		if err != nil {
			// un-recoverable error
			uc.logger.Errorw("could not transform request to event", "request", request.String())
			continue
		}
		events[key] = event
	}

	errs := uc.infra.Stream.Publisher("attempt_trigger_exec").Pub(ctx, events)
	for key := range events {
		if err, ok := errs[key]; ok {
			ko.Set(key, err)
			continue
		}

		ok.Append(key)
	}

	return ok, ko
}

// create attempts
func (uc *trigger) create(ctx context.Context, delay int64, requests []entities.Request) (*safe.Slice[string], *safe.Map[error]) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}
	attempts := []entities.Attempt{}

	now := uc.infra.Timer.Now()
	next := now.Add(time.Duration(delay) * time.Millisecond)

	for _, request := range requests {
		attempts = append(attempts, entities.Attempt{
			ReqId: request.Id,
			MsgId: request.MsgId,
			AppId: request.AppId,
			Tier:  request.Tier,

			ScheduledAt: now.UnixMilli(),
			Status:      0,

			ScheduleCounter: 0,
			ScheduleNext:    next.UnixMilli(),

			CompletedAt: 0,
		})
	}

	ids, err := uc.repositories.Datastore().Attempt().Create(ctx, attempts)
	if err != nil {
		for _, attempt := range attempts {
			ko.Set(attempt.ReqId, err)
		}
	}

	ok.Append(ids...)

	return ok, ko
}
