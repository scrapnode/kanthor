package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type WorkspaceCredentials interface {
	Authenticate(ctx context.Context, req *WorkspaceCredentialsAuthenticateIn) (*WorkspaceCredentialsAuthenticateOut, error)
}

type workspaceCredentials struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
