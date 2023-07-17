package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

type Workspace interface {
	Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error)
	ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error)
}

type WorkspaceGetReq struct {
	Id string
}

type WorkspaceGetRes struct {
	Workspace *entities.Workspace
}

type WorkspaceGetByAccountReq struct {
	WorkspaceId string
	Account     *authenticator.Account
}

type WorkspaceGetByAccountRes struct {
	Workspace *entities.Workspace
}

type WorkspaceListOfAccountReq struct {
	Account *authenticator.Account
}

type WorkspaceListOfAccountRes struct {
	Workspaces []entities.Workspace
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories
}
