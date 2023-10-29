package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc/pool"
)

type ForwarderSendReq struct {
	Concurrency int
	Requests    []entities.Request
}

func ValidateForwarderSendReqRequest(prefix string, item *entities.Request) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("request.id", item.Id, entities.IdNsReq),
		validator.StringStartsWith("request.msg_id", item.MsgId, entities.IdNsMsg),
		validator.StringStartsWith("request.ep_id", item.EpId, entities.IdNsEp),
		validator.StringRequired("request.tier", item.Tier),
		validator.StringStartsWith("request.app_id", item.AppId, entities.IdNsApp),
		validator.StringRequired("request.type", item.Type),
		validator.MapNotNil[string, string]("request.metadata", item.Metadata),
		validator.StringRequired("request.body", item.Body),
		validator.StringUri("request.uri", item.Uri),
		validator.StringRequired("request.method", item.Method),
	)
}

func (req *ForwarderSendReq) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("requests", req.Requests),
		validator.NumberGreaterThan("concurrency", req.Concurrency, 1),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Requests, func(i int, item *entities.Request) error {
			prefix := fmt.Sprintf("requests[%d]", i)
			return ValidateForwarderSendReqRequest(prefix, item)
		}),
	)
}

type ForwarderSendRes struct {
	Success []string
	Error   map[string]error
}

func (uc *forwarder) Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error) {
	responses := safe.Slice[*entities.Response]{}

	// we don't need to implement global timeout as we did with scheduler
	// because for each request, we already configured the sender timeout
	p := pool.New().WithMaxGoroutines(req.Concurrency)
	for _, r := range req.Requests {
		request := r
		p.Go(func() {
			response := uc.send(ctx, &request)
			responses.Append(response)
		})
	}
	p.Wait()

	events := map[string]*streaming.Event{}
	kv := responses.Data()
	for _, response := range kv {
		event, err := transformation.EventFromResponse(response)
		if err != nil {
			// un-recoverable error
			uc.logger.Errorw("could not transform response to event", "response", response.String())
			continue
		}

		events[response.ReqId] = event
	}

	ok := []string{}
	ko := map[string]error{}

	errs := uc.publisher.Pub(ctx, events)
	for reqId := range events {
		if err, ok := errs[reqId]; ok {
			ko[reqId] = err
			continue
		}
		ok = append(ok, reqId)
	}

	return &ForwarderSendRes{Success: ok, Error: ko}, nil
}

func (uc *forwarder) send(ctx context.Context, request *entities.Request) *entities.Response {
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
			if entities.Is5xx(res.Status) {
				return res, errors.New(http.StatusText(res.Status))
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
		response.Status = -1
	}

	if res != nil {
		response.Status = res.Status
		response.Uri = res.Uri
		response.Headers.FromHTTP(res.Headers)
		response.Body = string(res.Body)
	}

	return response
}
