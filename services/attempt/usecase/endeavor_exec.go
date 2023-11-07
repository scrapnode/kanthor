package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/safe"
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
		validator.NumberGreaterThan("concurrency", in.Concurrency, 1),
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
}

func (uc *endeavor) Exec(ctx context.Context, in *EndeavorExecIn) (*EndeavorExecOut, error) {
	responses := safe.Map[*entities.Response]{}

	requests, err := uc.requests(ctx, in)
	if err != nil {
		return nil, err
	}

	ok := safe.Slice[string]{}
	ko := safe.Map[error]{}

	// we don't need to implement global timeout as we did with scheduler
	// because for each request, we already configured the sender timeout
	p := pool.New().WithMaxGoroutines(in.Concurrency)
	for ref, r := range requests {
		refId := ref
		request := r
		p.Go(func() {
			response := uc.send(ctx, request)
			responses.Set(refId, response)

			if sender.Is5xxStatus(response.Status) {
				next := uc.infra.Timer.Now().Add(time.Millisecond * time.Duration(uc.conf.Endeavor.Executor.RescheduleDelay))
				err := uc.repositories.Datastore().Attempt().MarkReschedule(ctx, response.ReqId, next.UnixMilli())
				if err != nil {
					ko.Set(refId, err)
					return
				}

				ok.Append(refId)
				return
			}

			err := uc.repositories.Datastore().Attempt().MarkComplete(ctx, response.ReqId, response)
			if err != nil {
				ko.Set(refId, err)
				return
			}

			ok.Append(refId)

		})
	}
	p.Wait()

	events := map[string]*streaming.Event{}
	kv := responses.Data()
	for refId, response := range kv {
		event, err := transformation.EventFromResponse(response)
		if err != nil {
			// un-recoverable error
			uc.logger.Errorw("could not transform response to event", "response", response.String())
			continue
		}

		events[refId] = event
	}

	errs := uc.publisher.Pub(ctx, events)
	for refId, err := range errs {
		uc.logger.Errorw("unable to publish response event of request", "req_id", refId, "err", err.Error())
	}

	return &EndeavorExecOut{Success: ok.Data(), Error: ko.Data()}, nil
}

func (uc *endeavor) requests(ctx context.Context, in *EndeavorExecIn) (map[string]*entities.Request, error) {
	inIds := []string{}
	for _, attempt := range in.Attempts {
		inIds = append(inIds, attempt.ReqId)
	}
	requests, err := uc.repositories.Datastore().Request().ListByIds(ctx, inIds)
	if err != nil {
		return nil, err
	}

	returning := map[string]*entities.Request{}
	for _, request := range requests {
		returning[request.Id] = &request
	}
	return returning, nil
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
			if sender.Is5xxStatus(res.Status) {
				return res, errors.New(sender.StatusText(res.Status))
			}

			return res, nil
		},
		func(err error) error {
			return err
		},
	)

	response := &entities.Response{
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
	response.Metadata.Merge(request.Metadata)
	response.GenId()
	response.SetTS(uc.infra.Timer.Now())

	// IMPORTANT: we have an anti-pattern response that returns both error && response to trigger circuit breaker
	// so we should test both error and response seperately
	if err != nil {
		uc.logger.Errorw(err.Error(), "req_id", request.Id, "ep_id", request.EpId)
		response.Error = err.Error()
		response.Status = sender.Status(err.Error())
	}

	if res != nil {
		response.Status = res.Status
		response.Uri = res.Uri
		response.Headers.FromHTTP(res.Headers)
		response.Body = string(res.Body)
	}

	return response
}
