package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/status"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc/pool"
)

type EndeavorExecIn struct {
	Concurrency int

	Attempts map[string]*entities.Attempt
}

func ValidateEndeavorExecInAttempt(prefix string, item *entities.Attempt) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("attempt.req_id", item.ReqId, entities.IdNsReq),
		validator.StringStartsWith("attempt.app_id", item.AppId, entities.IdNsApp),
		validator.StringRequired("attempt.tier", item.Tier),
		validator.StringStartsWith("attempt.rest_id", item.ResId, entities.IdNsRes),
	)
}

func (in *EndeavorExecIn) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("concurrency", in.Concurrency, 0),
		validator.MapRequired("attempts", in.Attempts),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Attempts, func(refId string, item *entities.Attempt) error {
			prefix := fmt.Sprintf("attempts.%s", refId)
			return ValidateEndeavorExecInAttempt(prefix, item)
		}),
	)
}

type EndeavorExecOut struct {
	Success []string
	Error   map[string]error

	Rescheduled []string
	Completed   []string
}

func (uc *endeavor) Exec(ctx context.Context, in *EndeavorExecIn) (*EndeavorExecOut, error) {
	responses := safe.Map[*entities.Response]{}

	eventIdRefs := map[string]string{}
	for eventId, attempt := range in.Attempts {
		eventIdRefs[attempt.ReqId] = eventId
	}

	requests, err := uc.requests(ctx, in)
	if err != nil {
		return nil, err
	}

	ok := safe.Slice[string]{}
	ko := safe.Map[error]{}
	rescheduled := safe.Slice[string]{}
	completed := safe.Slice[string]{}

	// we don't need to implement global timeout as we did with scheduler
	// because for each request, we already configured the sender timeout
	p := pool.New().WithMaxGoroutines(in.Concurrency)
	for k, v := range requests {
		reqId := k
		request := v
		refId := eventIdRefs[reqId]
		p.Go(func() {
			response := uc.send(ctx, request)
			responses.Set(reqId, response)

			// reschedule for certaintly response type
			if response.Reschedulable() {
				next := uc.infra.Timer.Now().Add(time.Millisecond * time.Duration(uc.conf.Endeavor.Executor.RescheduleDelay))
				err := uc.repositories.Datastore().Attempt().MarkReschedule(ctx, response.ReqId, next.UnixMilli())
				if err != nil {
					ko.Set(refId, err)
					return
				}

				rescheduled.Append(reqId)
				ok.Append(refId)
				return
			}

			// otherwise mark the request as complete not matter what status it is (even though the status is fail)
			err := uc.repositories.Datastore().Attempt().MarkComplete(ctx, response.ReqId, response)
			if err != nil {
				ko.Set(reqId, err)
				return
			}

			completed.Append(reqId)
			ok.Append(refId)

		})
	}
	p.Wait()

	events := map[string]*streaming.Event{}
	kv := responses.Data()
	for reqId, response := range kv {
		event, err := transformation.EventFromResponse(response)
		if err != nil {
			// un-recoverable error, don't need to regconize it with KO map
			uc.logger.Errorw("could not transform response to event", "response", response.String())
			continue
		}

		events[reqId] = event
	}

	errs := uc.publisher.Pub(ctx, events)
	for reqId, err := range errs {
		refId := eventIdRefs[reqId]
		ko.Set(refId, err)

		uc.logger.Errorw("unable to publish response event of request", "req_id", reqId, "err", err.Error())
	}

	out := &EndeavorExecOut{
		Success:     ok.Data(),
		Error:       ko.Data(),
		Rescheduled: rescheduled.Data(),
		Completed:   completed.Data(),
	}
	return out, nil
}

func (uc *endeavor) requests(ctx context.Context, in *EndeavorExecIn) (map[string]*entities.Request, error) {
	maps := map[string]map[string][]string{}
	for _, attempt := range in.Attempts {
		if _, ok := maps[attempt.AppId]; !ok {
			maps[attempt.AppId] = map[string][]string{}
		}
		if _, ok := maps[attempt.AppId][attempt.MsgId]; !ok {
			maps[attempt.AppId][attempt.MsgId] = []string{}
		}
		maps[attempt.AppId][attempt.MsgId] = append(maps[attempt.AppId][attempt.MsgId], attempt.ReqId)
	}

	return uc.repositories.Datastore().Request().ListByIds(ctx, maps)
}

func (uc *endeavor) send(ctx context.Context, request *entities.Request) *entities.Response {
	res, err := circuitbreaker.Do[sender.Response](
		uc.infra.CircuitBreaker,
		request.EpId,
		func() (interface{}, error) {
			in := &sender.Request{
				Method:  request.Method,
				Headers: request.Headers.ToHTTP(),
				Uri:     request.Uri,
				Body:    []byte(request.Body),
			}

			res, err := uc.infra.Send(ctx, in)
			if err != nil {
				return nil, err
			}

			// sending is success, but we got remote server error
			// must use custom error here to trigger circuit breaker
			if status.Is5xx(res.Status) {
				return res, errors.New(status.Text(res.Status))
			}

			return res, nil
		},
		func(err error) error {
			return err
		},
	)

	doc := &entities.Response{
		MsgId:    request.MsgId,
		EpId:     request.EpId,
		ReqId:    request.Id,
		Tier:     request.Tier,
		AppId:    request.AppId,
		Type:     request.Type,
		Headers:  entities.Header{},
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	doc.Metadata.Merge(request.Metadata)
	doc.Id = suid.New(entities.IdNsRes)
	doc.SetTS(uc.infra.Timer.Now())

	// IMPORTANT: we have an anti-pattern response that returns both error && response to trigger circuit breaker
	// so we should test both error and response seperately
	if err != nil {
		uc.logger.Errorw(err.Error(), "req_id", request.Id, "ep_id", request.EpId)
		doc.Error = err.Error()
		doc.Status = status.Code(err.Error())
	}

	if res != nil {
		doc.Status = res.Status
		doc.Uri = res.Uri
		doc.Headers.FromHTTP(res.Headers)
		doc.Body = string(res.Body)
	}

	return doc
}
