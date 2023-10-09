package dispatcher

import (
	"context"
	"errors"
	"net/http"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

type ForwarderSendReq struct {
	Request ForwarderSendReqRequest
}

func (req *ForwarderSendReq) Validate() error {
	if err := req.Request.Validate(); err != nil {
		return err
	}
	return nil
}

type ForwarderSendReqRequest struct {
	Id    string
	MsgId string
	EpId  string

	Tier     string
	AppId    string
	Type     string
	Metadata entities.Metadata

	Headers entities.Header
	Body    []byte
	Uri     string
	Method  string
}

func (req *ForwarderSendReqRequest) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,

		validator.StringStartsWith("request.id", req.Id, entities.IdNsReq),
		validator.StringStartsWith("request.msg_id", req.MsgId, entities.IdNsMsg),
		validator.StringStartsWith("request.ep_id", req.EpId, entities.IdNsEp),
		validator.StringRequired("request.tier", req.Tier),
		validator.StringStartsWith("request.app_id", req.AppId, entities.IdNsApp),
		validator.StringRequired("request.type", req.Type),
		validator.MapNotNil[string, string]("request.metadata", req.Metadata),
		validator.SliceRequired("request.body", req.Body),
		validator.StringUri("request.uri", req.Uri),
		validator.StringRequired("request.method", req.Method),
	)
}

type ForwarderSendRes struct {
	Response entities.Response
}

func (uc *forwarder) Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error) {
	request := &sender.Request{
		Method:  req.Request.Method,
		Headers: req.Request.Headers.Header,
		Uri:     req.Request.Uri,
		Body:    req.Request.Body,
	}

	// @TODO: apply rate limit to endpoint
	response, err := circuitbreaker.Do[sender.Response](
		uc.cb,
		req.Request.EpId,
		func() (interface{}, error) {
			res, err := uc.dispatch(request)
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

	res := entities.Response{
		MsgId:    req.Request.MsgId,
		EpId:     req.Request.EpId,
		ReqId:    req.Request.Id,
		Tier:     req.Request.Tier,
		AppId:    req.Request.AppId,
		Type:     req.Request.Type,
		Headers:  entities.NewHeader(),
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	res.Metadata.Merge(req.Request.Metadata)
	res.GenId()
	res.SetTS(uc.timer.Now())

	// IMPORTANT: we have an anti-pattern case that returns both error && response to trigger circuit breaker
	// so we should test both error and response seperately
	if err != nil {
		uc.logger.Errorw(err.Error(), "req_id", req.Request.Id, "ep_id", req.Request.EpId)
		res.Error = err.Error()
		res.Status = -1
	}

	if response != nil {
		res.Status = response.Status
		res.Uri = response.Uri
		res.Headers.Merge(entities.Header{Header: response.Headers})
		res.Body = response.Body
	}

	event, err := transformation.EventFromResponse(&res)
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Pub(ctx, event); err != nil {
		return nil, err
	}

	return &ForwarderSendRes{Response: res}, nil
}
