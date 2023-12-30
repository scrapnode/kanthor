package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Endpoint interface {
	ListMessage(ctx context.Context, in *EndpointListMessageIn) (*EndpointListMessageOut, error)
	GetMessage(ctx context.Context, in *EndpointGetMessageIn) (*EndpointGetMessageOut, error)
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
