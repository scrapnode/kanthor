package dispatcher

import (
	"context"
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
	AttId string

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

		validator.StringStartsWith("request.id", req.Id, "req_"),
		validator.StringStartsWith("request.att_id", req.AttId, "att_"),
		validator.StringRequired("request.tier", req.Tier),
		validator.StringStartsWith("request.app_id", req.AppId, "app_"),
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
		req.Request.Metadata.Get(entities.MetaEpId),
		func() (interface{}, error) {
			return uc.dispatch(request)
		},
		func(err error) error {
			return err
		},
	)

	res := entities.Response{
		AttId:    req.Request.AttId,
		Tier:     req.Request.Tier,
		AppId:    req.Request.AppId,
		Type:     req.Request.Type,
		Headers:  entities.Header{Header: http.Header{}},
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	res.Metadata.Merge(req.Request.Metadata)
	res.GenId()
	res.SetTS(uc.timer.Now(), uc.conf.Bucket.Layout)

	// either error was happened or not, we need to publish response event, so we can handle custom logic later
	// example use cases are retry, notification, i.e
	if err == nil {
		res.Status = response.Status
		res.Uri = response.Uri
		res.Headers.Merge(entities.Header{Header: response.Headers})
		res.Body = response.Body
	} else {
		uc.logger.Errorw(err.Error(), "ep_id", req.Request.Metadata.Get(entities.MetaEpId), "req_id", req.Request.Id)
		res.Status = -1
		res.Error = err.Error()
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
