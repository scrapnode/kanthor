package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Workspace interface {
	Create(ctx context.Context, in *WorkspaceCreateIn) (*WorkspaceCreateOut, error)
	Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error)
	List(ctx context.Context, in *WorkspaceListIn) (*WorkspaceListOut, error)
	Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error)

	Setup(ctx context.Context, in *WorkspaceSetupIn) (*WorkspaceSetupOut, error)
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
