package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Application interface {
	ListMessage(ctx context.Context, in *ApplicationListMessageIn) (*ApplicationListMessageOut, error)
	GetMessage(ctx context.Context, in *ApplicationGetMessageIn) (*ApplicationGetMessageOut, error)
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
