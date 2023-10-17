package portal

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Workspace interface {
	Setup(ctx context.Context, req *WorkspaceSetupReq) (*WorkspaceSetupRes, error)

	Update(ctx context.Context, req *WorkspaceUpdateReq) (*WorkspaceUpdateRes, error)
	Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error)
}

type workspace struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
