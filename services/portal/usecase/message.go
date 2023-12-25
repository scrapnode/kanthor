package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Message interface {
	List(ctx context.Context, in *MessageListIn) (*MessageListOut, error)
	Get(ctx context.Context, in *MessageGetIn) (*MessageGetOut, error)
}

type message struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
