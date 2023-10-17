package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type message struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
	repos     repos.Repositories
}
