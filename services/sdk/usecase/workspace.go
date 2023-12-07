package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type Workspace interface {
	Authenticate(ctx context.Context, req *WorkspaceAuthenticateIn) (*WorkspaceAuthenticateOut, error)
	Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error)
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
