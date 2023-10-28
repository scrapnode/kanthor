package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/planner"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
)

type TriggerExecReq struct {
	Size    int
	Timeout int64

	AttemptDelay int64
	Triggers     map[string]*entities.AttemptTrigger
}

func (req *TriggerExecReq) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("size", req.Size, 1),
		validator.NumberGreaterThan("timeout", int(req.Timeout), 1000),
		validator.NumberGreaterThan("attempt_delay", req.AttemptDelay, 60000),
		validator.MapRequired("triggers", req.Triggers),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Map(req.Triggers, func(key string, item *entities.AttemptTrigger) error {
			prefix := fmt.Sprintf("triggers.%s", key)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".app.id", req.Triggers[key].AppId, entities.IdNsApp),
				validator.StringRequired(prefix+".tier", req.Triggers[key].Tier),
				validator.NumberGreaterThan(prefix+".from", req.Triggers[key].From, 0),
				validator.NumberGreaterThan(prefix+".to", int(req.Triggers[key].To), 0),
			)
		}),
	)
}

type TriggerExecRes struct {
	Success []string
	Error   map[string]error
}

func (uc *trigger) Exec(ctx context.Context, req *TriggerExecReq) (*TriggerExecRes, error) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}

	// timeout duration will be scaled based on how many triggers you have
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(req.Timeout))
	defer cancel()

	var wg conc.WaitGroup
	for _, item := range req.Triggers {
		trigger := item
		wg.Go(func() {
			resp, err := uc.consume(ctx, trigger, req.Size, req.AttemptDelay)
			if err != nil {
				uc.logger.Errorw("unable to consume attempt trigger", "trigger", trigger.String(), "err", err.Error())
				return
			}

			ko.Merge(resp.Error)
			ok.Append(resp.Success...)
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
		return &TriggerExecRes{Success: ok.Data(), Error: ko.Data()}, nil
	case <-timeout.Done():
		// actually we may have some success entries, but we can ignore them
		// let cronjob pickup and retry them redundantly
		return nil, ctx.Err()
	}
}

func (uc *trigger) consume(
	ctx context.Context,
	trigger *entities.AttemptTrigger,
	size int,
	delay int64,
) (*TriggerExecRes, error) {
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
	msgIds, requests, err := uc.examine(ctx, trigger.AppId, applicable, from, to)
	if err != nil {
		return nil, err
	}

	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}

	uc.schedule(ctx, ok, ko, size, applicable, msgIds)
	uc.create(ctx, ok, ko, size, delay, requests)

	res := &TriggerExecRes{Success: ok.Data(), Error: ko.Data()}
	return res, nil
}

// examine messages, requests and responses
func (uc *trigger) examine(
	ctx context.Context,
	appId string,
	applicable *planner.Applicable,
	from, to time.Time,
) ([]string, []repositories.Req, error) {
	messages, msgIds, err := uc.scan(ctx, appId, from, to)
	if err != nil {
		return nil, nil, err
	}

	requests, err := uc.repositories.Request().Scan(ctx, appId, msgIds, from, to)
	if err != nil {
		return nil, nil, err
	}

	responses, err := uc.repositories.Response().Scan(ctx, appId, msgIds, from, to)
	if err != nil {
		return nil, nil, err
	}

	schedulable := []string{}
	attemptable := []repositories.Req{}

	status := uc.hash(requests, responses)
	for _, message := range messages {
		for _, ep := range applicable.EndpointMap {
			reqId, hasReq := status[reqkey(message.Id, ep.Id)]
			if !hasReq {
				// no request -> must schedule message again -> don't create any attempt
				schedulable = append(schedulable, message.Id)
				continue
			}

			_, hasRes := status[reskey(message.Id, ep.Id)]
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

func (uc *trigger) scan(ctx context.Context, appId string, from, to time.Time) (map[string]repositories.Msg, []string, error) {
	cursor, err := uc.infra.Cache.StringGet(ctx, "kanthor.usecases.attempt.message.scan")
	if err != nil && !errors.Is(err, cache.ErrEntryNotFound) {
		return nil, nil, err
	}

	messages, err := uc.repositories.Message().Scan(ctx, appId, from, to, cursor)
	if err != nil {
		return nil, nil, err
	}

	if len(messages) > 0 {
		cursor = messages[len(messages)-1].Id
	}

	err = uc.infra.Cache.StringSet(ctx, "kanthor.usecases.attempt.message.scan", cursor, time.Hour)
	if err != nil {
		uc.logger.Errorw("unable to set scan cursor to reuse later", "err", err.Error(), "cursor", cursor)
	}

	ids := []string{}
	returning := map[string]repositories.Msg{}
	for _, msg := range messages {
		returning[msg.Id] = msg
		ids = append(ids, msg.Id)
	}

	return returning, ids, nil
}

func (uc *trigger) applicable(ctx context.Context, appId string) (*planner.Applicable, error) {
	key := utils.Key("scheduler", appId)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour, func() (*planner.Applicable, error) {
		endpoints, err := uc.repositories.Endpoint().List(ctx, appId)
		if err != nil {
			return nil, err
		}
		returning := &planner.Applicable{EndpointMap: map[string]entities.Endpoint{}}
		for _, ep := range endpoints {
			returning.EndpointMap[ep.Id] = ep
		}

		rules, err := uc.repositories.Endpoint().Rules(ctx, appId)
		if err != nil {
			return nil, err
		}
		returning.Rules = rules

		return returning, nil
	})
}

func (uc *trigger) hash(requests map[string]repositories.Req, responses map[string]repositories.Res) map[string]string {
	returning := map[string]string{}

	for _, request := range requests {
		// for checking whether we have scheduled a request for an endpoint or not
		// if no request was scheduled, we should schedule it instead of create an attempt
		key := reqkey(request.MsgId, request.EpId)
		returning[key] = request.Id
	}

	for _, response := range responses {
		key := reskey(response.MsgId, response.EpId)

		// we already recognized that the endpoint had success status, don't need to check any more
		if _, has := returning[key]; has {
			continue
		}

		// status is ok, saved the success response id
		if entities.Is2xx(response.Status) {
			returning[key] = response.Id
		}
	}

	return returning
}

func reqkey(msgId, epId string) string {
	return fmt.Sprintf("%s/%s/req", msgId, epId)
}

func reskey(msgId, epId string) string {
	return fmt.Sprintf("%s/%s/res", msgId, epId)
}

// schedule messages
func (uc *trigger) schedule(
	ctx context.Context,
	ok *safe.Slice[string],
	ko *safe.Map[error],
	size int,
	applicable *planner.Applicable,
	msgIds []string,
) {
	requests := []entities.Request{}

	for i := 0; i < len(msgIds); i += size {
		j := utils.ChunkNext(i, len(msgIds), size)

		messages, err := uc.repositories.Message().ListByIds(ctx, msgIds[i:j])
		if err != nil {
			for _, id := range msgIds[i:j] {
				ko.Set(id, err)
			}
			continue
		}

		for _, message := range messages {
			reqs, logs := planner.Requests(&message, applicable, uc.infra.Timer)
			if len(logs) > 0 {
				for _, l := range logs {
					uc.logger.Warnw(l[0].(string), l[1:]...)
				}
			}
			requests = append(requests, reqs...)
		}
	}

	if len(requests) == 0 {
		return
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
}

// create attempts
func (uc *trigger) create(
	ctx context.Context,
	ok *safe.Slice[string],
	ko *safe.Map[error],
	size int,
	delay int64,
	requests []repositories.Req,
) {
	attempts := []entities.Attempt{}

	now := uc.infra.Timer.Now()
	next := now.Add(time.Duration(delay) * time.Millisecond)

	for _, request := range requests {
		attempts = append(attempts, entities.Attempt{
			ReqId:        request.Id,
			Tier:         request.Tier,
			ScheduleNext: next.UnixMilli(),
			ScheduledAt:  now.UnixMilli(),
		})
	}

	for i := 0; i < len(attempts); i += size {
		j := utils.ChunkNext(i, len(attempts), size)

		ids, err := uc.repositories.Attempt().Create(ctx, attempts[i:j])
		if err != nil {
			for _, attempt := range attempts[i:j] {
				ko.Set(attempt.ReqId, err)
			}
			continue
		}

		ok.Append(ids...)
	}
}
