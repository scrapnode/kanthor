package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/assessor"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/status"
	"github.com/scrapnode/kanthor/internal/domain/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/attempt/repositories/ds"
)

func (uc *trigger) Perform(
	ctx context.Context,
	appId string,
	msgs map[string]entities.Message,
	applicable *assessor.Assets,
	attemptDelay int64,
) (*TriggerExecOut, error) {
	messages, requests, err := uc.examine(ctx, appId, applicable, msgs)
	if err != nil {
		return nil, err
	}

	out := &TriggerExecOut{
		Success:   []string{},
		Error:     map[string]error{},
		Scheduled: []string{},
		Created:   []string{},
	}

	scheduledok, scheduledko := uc.schedule(ctx, messages, applicable)
	out.Success = append(out.Success, scheduledok.Data()...)
	for k, v := range scheduledko.Data() {
		out.Error[k] = v
	}
	out.Scheduled = append(out.Scheduled, scheduledok.Data()...)

	createdok, createdko := uc.create(ctx, requests, attemptDelay)
	out.Success = append(out.Success, createdok.Data()...)
	for k, v := range createdko.Data() {
		out.Error[k] = v
	}
	out.Created = append(out.Created, createdok.Data()...)

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
func (uc *trigger) schedule(ctx context.Context, messages []entities.Message, applicable *assessor.Assets) (*safe.Slice[string], *safe.Map[error]) {
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
func (uc *trigger) create(ctx context.Context, requests []entities.Request, delay int64) (*safe.Slice[string], *safe.Map[error]) {
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
