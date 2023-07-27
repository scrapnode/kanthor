package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Workspace interface {
	Import(ctx context.Context, req *WorkspaceImportReq) (*WorkspaceImportRes, error)
}

type WorkspaceImportReq struct {
	Workspaces           []entities.Workspace            `validate:"required"`
	WorkspaceTiers       []entities.WorkspaceTier        `validate:"required"`
	WorkspaceCredentials []entities.WorkspaceCredentials `validate:"required"`
	Applications         []entities.Application          `validate:"required"`
	Endpoints            []entities.Endpoint             `validate:"required"`
	EndpointRules        []entities.EndpointRule         `validate:"required"`
}

type WorkspaceImportRes struct {
	WorkspaceIds            []string
	WorkspaceTierIds        []string
	WorkspaceCredentialsIds []string
	ApplicationIds          []string
	EndpointIds             []string
	EndpointRuleIds         []string
}

type workspace struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
