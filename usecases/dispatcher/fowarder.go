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
)

type Forwarder interface {
	Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error)
}

type ForwarderSendReq struct {
	Request ForwarderSendReqRequest `validate:"required"`
}

type ForwarderSendReqRequest struct {
	Id    string `validate:"required,startswith=req_"`
	AttId string `validate:"required,startswith=att_"`

	Tier     string            `validate:"required"`
	AppId    string            `validate:"required,startswith=app_"`
	Type     string            `validate:"required"`
	Metadata entities.Metadata `validate:"required"`

	Headers entities.Header `validate:"required"`
	Body    []byte          `validate:"required"`
	Uri     string          `validate:"required,uri"`
	Method  string          `validate:"required"`
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
