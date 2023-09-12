package scheduler

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/signature"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

type Request interface {
	Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error)
}

type RequestArrangeReq struct {
	Message RequestArrangeReqMessage `validate:"required"`
}

type RequestArrangeReqMessage struct {
	Id    string `validate:"required,startswith=msg_"`
	AttId string `validate:"required,startswith=att_"`

	Tier     string            `validate:"required"`
	AppId    string            `validate:"required,startswith=app_"`
	Type     string            `validate:"required"`
	Metadata entities.Metadata `validate:"required"`
	Headers  entities.Header   `validate:"required"`
	Body     []byte            `validate:"required"`
}

type RequestArrangeRes struct {
	Entities    []structure.BulkRes[entities.Request]
	FailKeys    []string
	SuccessKeys []string
}

type request struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	signature signature.Signature
	publisher streaming.Publisher
	cache     cache.Cache
	metrics   metric.Metrics
	repos     repos.Repositories
}
