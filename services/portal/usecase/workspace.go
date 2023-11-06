package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Workspace interface {
	Setup(ctx context.Context, in *WorkspaceSetupIn) (*WorkspaceSetupOut, error)

	Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error)
	Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error)
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
