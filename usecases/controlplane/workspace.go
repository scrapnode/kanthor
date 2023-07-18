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
	Id string `json:"id" validate:"required"`
}

type WorkspaceGetRes struct {
	Workspace *entities.Workspace `json:"workspace"`
}

type WorkspaceGetByAccountReq struct {
	Account     *authenticator.Account `json:"account" validate:"required"`
	WorkspaceId string                 `json:"workspace_id" validate:"required"`
}

type WorkspaceGetByAccountRes struct {
	Workspace *entities.Workspace `json:"workspace"`
}

type WorkspaceListOfAccountReq struct {
	Account *authenticator.Account `json:"account" validate:"required"`
}

type WorkspaceListOfAccountRes struct {
	Workspaces []entities.Workspace `json:"workspaces"`
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
