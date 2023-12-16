package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/assessor"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/status"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services/attempt/repositories/ds"
	"github.com/sourcegraph/conc/pool"
)

type TriggerPerformIn struct {
	AppId        string
	Concurrency  int
	ArrangeDelay int64

	Applicable *assessor.Assets
	Messages   map[string]*entities.Message
}

func (in *TriggerPerformIn) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.NumberGreaterThan("concurrency", in.Concurrency, 0),
		validator.NumberGreaterThan("arrange_delay", in.ArrangeDelay, 60000),
		validator.PointerNotNil("applicable", in.Applicable),
		validator.MapRequired("messages", in.Messages),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Messages, func(id string, item *entities.Message) error {
			prefix := fmt.Sprintf("messages.%s", id)
			return ValidateWarehousePutInMessage(prefix, item)
		}),
	)
}

func ValidateWarehousePutInMessage(prefix string, message *entities.Message) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", message.Id, entities.IdNsMsg),
		validator.NumberGreaterThan(prefix+".timestamp", message.Timestamp, 0),
		validator.StringRequired(prefix+".tier", message.Tier),
		validator.StringStartsWith(prefix+".app_id", message.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", message.Type),
		validator.StringRequired(prefix+".body", message.Body),
	)
}

func (uc *trigger) Perform(ctx context.Context, in *TriggerPerformIn) (*TriggerExecOut, error) {
	messages, requests, err := uc.examine(ctx, in.AppId, in.Applicable, in.Messages)
	if err != nil {
		return nil, err
	}

	out := &TriggerExecOut{
		Success:   []string{},
		Error:     map[string]error{},
		Scheduled: []string{},
		Created:   []string{},
	}

	scheduledok, scheduledko := uc.schedule(ctx, messages, in)
	out.Success = append(out.Success, scheduledok.Data()...)
	for k, v := range scheduledko.Data() {
		out.Error[k] = v
	}
	out.Scheduled = append(out.Scheduled, scheduledok.Data()...)

	createdok, createdko := uc.create(ctx, requests, in)
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
	messages map[string]*entities.Message,
) (map[string]*entities.Message, map[string]*entities.Request, error) {
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

	schedulable := map[string]*entities.Message{}
	attemptable := map[string]*entities.Request{}

	status := uc.hash(requests, responses)
	for _, message := range messages {
		for _, ep := range applicable.EndpointMap {
			reqId, hasReq := status[ds.ReqKey(message.Id, ep.Id)]
			if !hasReq {
				// no request -> must schedule message again -> don't create any attempt
				schedulable[message.Id] = message
				continue
			}

			_, hasRes := status[ds.ResKey(message.Id, ep.Id)]
			if !hasRes {
				// has request + no success response -> create an attempt
				attemptable[reqId] = requests[reqId]
				continue
			}

			// has success response, ignore
		}
	}

	return schedulable, attemptable, nil
}

func (uc *trigger) hash(requests map[string]*entities.Request, responses map[string]*ds.ResponseStatusRow) map[string]string {
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
func (uc *trigger) schedule(ctx context.Context, messages map[string]*entities.Message, in *TriggerPerformIn) (*safe.Slice[string], *safe.Map[error]) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}
	requests := map[string]*entities.Request{}

	for _, message := range messages {
		reqs, logs := assessor.Requests(message, in.Applicable, uc.infra.Timer)
		if len(logs) > 0 {
			for _, l := range logs {
				uc.logger.Warnw(l[0].(string), l[1:]...)
			}
		}
		for reqId, req := range reqs {
			requests[reqId] = req
		}
	}

	if len(requests) == 0 {
		return ok, ko
	}

	events := map[string]*streaming.Event{}
	for _, request := range requests {
		key := utils.Key(request.AppId, request.MsgId, request.EpId, request.Id)
		event, err := transformation.EventFromRequest(request)
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
func (uc *trigger) create(ctx context.Context, requests map[string]*entities.Request, in *TriggerPerformIn) (*safe.Slice[string], *safe.Map[error]) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}
	attempts := []*entities.Attempt{}

	now := uc.infra.Timer.Now()
	next := now.Add(time.Duration(in.ArrangeDelay) * time.Millisecond)

	for _, request := range requests {
		attempts = append(attempts, &entities.Attempt{
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

	// hardcode the go routine to 1 because we are expecting stable throughput of database inserting
	p := pool.New().WithMaxGoroutines(1)
	for i := 0; i < len(attempts); i += in.Concurrency {
		j := utils.ChunkNext(i, len(attempts), in.Concurrency)

		items := attempts[i:j]
		p.Go(func() {
			reqIds, err := uc.repositories.Datastore().Attempt().Create(ctx, items)
			if err != nil {
				for _, attempt := range attempts {
					ko.Set(attempt.ReqId, err)
				}
			}

			ok.Append(reqIds...)
		})
	}
	p.Wait()

	return ok, ko
}
