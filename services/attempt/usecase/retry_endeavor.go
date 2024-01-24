package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/status"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc/pool"
)

type RetryEndeavorIn struct {
	Concurrency int
	Attempts    map[string]*entities.Attempt
}

func ValidateRetryEndeavorInAttempts(prefix string, attempt *entities.Attempt) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".req_id", attempt.ReqId, entities.IdNsReq),
		validator.StringStartsWith(prefix+".msg_id", attempt.MsgId, entities.IdNsMsg),
		validator.StringStartsWith(prefix+".app_id", attempt.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".tier", attempt.Tier),
	)
}

func (in *RetryEndeavorIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.MapRequired("attempts", in.Attempts),
		validator.Map(in.Attempts, func(refId string, item *entities.Attempt) error {
			prefix := fmt.Sprintf("attempts.%s", refId)
			return ValidateRetryEndeavorInAttempts(prefix, item)
		}),
	)
}

type RetryEndeavorOut struct {
	Success []string
	Error   map[string]error
}

func (uc *retry) Endeavor(ctx context.Context, in *RetryEndeavorIn) (*RetryEndeavorOut, error) {
	requests, err := uc.repositories.Datastore().Attempt().ListRequests(ctx, in.Attempts)
	if err != nil {
		return nil, err
	}

	responses := safe.Map[*entities.Response]{}
	// we don't need to implement global timeout as we did with scheduler
	// because for each request, we already configured the sender timeout
	p := pool.New().WithMaxGoroutines(int(in.Concurrency))
	for ref, r := range requests {
		refId := ref
		request := r
		p.Go(func() {
			response := uc.send(ctx, request)
			responses.Set(refId, response)
		})
	}
	p.Wait()

	events := map[string]*streaming.Event{}
	kv := responses.Data()
	for refId, response := range kv {
		event, err := transformation.EventFromResponse(response)
		if err != nil {
			// un-recoverable error
			uc.logger.Errorw("ATTEMPT.USECASE.RETRY.ENDEAVOR.EVENT_TRANSFORM.ERROR", "response", response.String())
			continue
		}

		events[refId] = event
	}

	ok := []string{}
	ko := map[string]error{}

	errs := uc.publisher.Pub(ctx, events)
	for refId := range events {
		if err, ok := errs[refId]; ok {
			ko[refId] = err
			continue
		}
		ok = append(ok, refId)
	}

	return &RetryEndeavorOut{Success: ok, Error: ko}, nil
}

func (uc *retry) send(ctx context.Context, request *entities.Request) *entities.Response {
	res, err := circuitbreaker.Do[sender.Response](
		uc.infra.CircuitBreaker,
		request.EpId,
		func() (interface{}, error) {
			req := &sender.Request{
				Method:  request.Method,
				Headers: request.Headers.ToHTTP(),
				Uri:     request.Uri,
				Body:    []byte(request.Body),
			}

			res, err := uc.infra.Send(ctx, req)
			if err != nil {
				return nil, err
			}

			// sending is success, but we got remote server error
			// must use custom error here to trigger circuit breaker
			if status.IsKO(res.Status) {
				return res, errors.New(http.StatusText(res.Status))
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
	doc.Id = identifier.New(entities.IdNsRes)
	doc.SetTS(uc.infra.Timer.Now())

	// IMPORTANT: we have an anti-pattern response that returns both error && response to trigger circuit breaker
	// so we should test both error and response seperately
	if err != nil {
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
