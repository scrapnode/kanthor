package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

type Workspace interface {
	ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error)
	GetByAccount(ctx context.Context, req *WorkspaceGetByAccountReq) (*WorkspaceGetByAccountRes, error)
}

type WorkspaceGetByAccountReq struct {
	WorkspaceId string
	Account     *authenticator.Account
}

type WorkspaceGetByAccountRes struct {
	Workspace *entities.Workspace
}

type WorkspaceListOfAccountReq struct {
	structure.ListReq
	Account *authenticator.Account
}

type WorkspaceListOfAccountRes struct {
	Workspaces []entities.Workspace
}

type worksapce struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
