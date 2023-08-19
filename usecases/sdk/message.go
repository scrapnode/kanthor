package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
	"net/http"
)

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type MessagePutReq struct {
	AppId string `validate:"required,startswith=app_"`
	Type  string `validate:"required"`

	Body     []byte `validate:"required"`
	Headers  http.Header
	Metadata entities.Metadata
}

type MessagePutRes struct {
	Msg *entities.Message
}

type message struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metrics.Metrics
	timer        timer.Timer
	cache        cache.Cache
	publisher    streaming.Publisher
	repos        repos.Repositories
}
