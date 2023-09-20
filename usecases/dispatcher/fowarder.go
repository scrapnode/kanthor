package dispatcher

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Forwarder interface {
	Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error)
}

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

type forwarder struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	dispatch  sender.Send
	cache     cache.Cache
	cb        circuitbreaker.CircuitBreaker
	metrics   metric.Metrics
}
