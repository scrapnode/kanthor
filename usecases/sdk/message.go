package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
	"net/http"
)

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type MessagePutReq struct {
	AppId string `json:"app_id" validate:"required,startswith=app_"`
	Type  string `json:"type" validate:"required"`

	Body     []byte            `json:"body" validate:"required"`
	Headers  http.Header       `json:"headers"`
	Metadata map[string]string `json:"metadata"`
}

type MessagePutRes struct {
	Msg *entities.Message
}

type message struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	publisher    streaming.Publisher
	repos        repos.Repositories
}
