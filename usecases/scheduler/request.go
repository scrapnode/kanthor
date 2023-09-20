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
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

type Request interface {
	Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error)
}

type RequestArrangeReq struct {
	Message RequestArrangeReqMessage
}

func (req *RequestArrangeReq) Validate() error {
	if err := req.Message.Validate(); err != nil {
		return err
	}
	return nil
}

type RequestArrangeReqMessage struct {
	Id    string
	AttId string

	Tier     string
	AppId    string
	Type     string
	Metadata entities.Metadata
	Headers  entities.Header
	Body     []byte
}

func (req *RequestArrangeReqMessage) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,

		validator.StringStartsWith("message.id", req.Id, "msg_"),
		validator.StringStartsWith("message.att_id", req.AttId, "att_"),
		validator.StringRequired("message.tier", req.Tier),
		validator.StringStartsWith("message.app_id", req.AppId, "app_"),
		validator.StringRequired("message.type", req.Type),
		validator.MapNotNil[string, string]("message.metadata", req.Metadata),
		validator.SliceRequired("message.body", req.Body),
	)
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
