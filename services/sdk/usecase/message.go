package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type Message interface {
	Put(ctx context.Context, in *MessagePutIn) (*MessagePutOut, error)
}

type message struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
}
