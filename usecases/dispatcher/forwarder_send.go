package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/pool"
)

type ForwarderSendReq struct {
	ChunkTimeout int64
	ChunkSize    int
	Requests     []entities.Request
}

func ValidateForwarderSendRequest(prefix string, item entities.Request) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("request.id", item.Id, entities.IdNsReq),
		validator.StringStartsWith("request.msg_id", item.MsgId, entities.IdNsMsg),
		validator.StringStartsWith("request.ep_id", item.EpId, entities.IdNsEp),
		validator.StringRequired("request.tier", item.Tier),
		validator.StringStartsWith("request.app_id", item.AppId, entities.IdNsApp),
		validator.StringRequired("request.type", item.Type),
		validator.MapNotNil[string, string]("request.metadata", item.Metadata),
		validator.SliceRequired("request.body", item.Body),
		validator.StringUri("request.uri", item.Uri),
		validator.StringRequired("request.method", item.Method),
	)
}

func (req *ForwarderSendReq) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("requests", req.Requests),
		validator.NumberGreaterThan("chunk_timeout", req.ChunkTimeout, 1000),
		validator.NumberGreaterThan("chunk_size", req.ChunkSize, 1),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Requests, func(i int, item entities.Request) error {
			prefix := fmt.Sprintf("requests[%d]", i)
			return ValidateForwarderSendRequest(prefix, item)
		}),
	)
}

type ForwarderSendRes struct {
	Success []string
	Error   map[string]error
}

func (uc *forwarder) Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error) {
	responses := []entities.Response{}

	// we don't need to implement global timeout as we did with scheduler
	// because for each request, we already configured the sender timeout
	var wg conc.WaitGroup
	for _, r := range req.Requests {
		request := r
		wg.Go(func() {
			response := uc.send(ctx, &request)
			responses = append(responses, *response)
		})
	}

	ok := &safe.Map[string]{}
	ko := &safe.Map[error]{}

	// but publishing need implementing the global timeout
	// timeout duration will be scaled based on how many responses you have
	duration := time.Duration(req.ChunkTimeout * int64(len(responses)+1))
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*duration)
	defer cancel()

	p := pool.New().WithMaxGoroutines(req.ChunkSize)
	for _, rs := range responses {
		response := rs
		p.Go(func() {

			event, err := transformation.EventFromResponse(&response)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform response to event", "response", response.String())
			}

			if err := uc.publisher.Pub(ctx, event); err != nil {
				ko.Set(response.ReqId, err)
				return
			}

			key := utils.Key(response.AppId, response.MsgId, response.EpId, response.ReqId, response.Id)
			ok.Set(key, response.ReqId)
		})
	}

	c := make(chan bool)
	defer close(c)

	go func() {
		p.Wait()
		c <- true
	}()

	select {
	case <-c:
		return &ForwarderSendRes{Success: ok.Keys(), Error: ko.Data()}, nil
	case <-timeout.Done():
		// context deadline exceeded, should set that error to remain requests
		for _, request := range req.Requests {
			if _, success := ok.Get(request.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(request.Id); !has {
				ko.Set(request.Id, ctx.Err())
			}
		}
		return &ForwarderSendRes{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *forwarder) send(ctx context.Context, request *entities.Request) *entities.Response {
	req := &sender.Request{
		Method:  request.Method,
		Headers: request.Headers.Header,
		Uri:     request.Uri,
		Body:    request.Body,
	}

	res, err := circuitbreaker.Do[sender.Response](
		uc.cb,
		request.EpId,
		func() (interface{}, error) {
			res, err := uc.dispatch(ctx, req)
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
		Headers:  entities.NewHeader(),
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	response.Metadata.Merge(request.Metadata)
	response.GenId()
	response.SetTS(uc.timer.Now())

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
		response.Headers.Merge(entities.Header{Header: res.Headers})
		response.Body = res.Body
	}

	return response
}
