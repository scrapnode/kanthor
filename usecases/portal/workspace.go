package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Workspace interface {
	Import(ctx context.Context, req *WorkspaceImportReq) (*WorkspaceImportRes, error)
	Export(ctx context.Context, req *WorkspaceExportReq) (*WorkspaceExportRes, error)

	Update(ctx context.Context, req *WorkspaceUpdateReq) (*WorkspaceUpdateRes, error)
	Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error)
}

type WorkspaceImportReq struct {
	Workspaces    []entities.Workspace    `validate:"required,gt=0"`
	Applications  []entities.Application  `validate:"required"`
	Endpoints     []entities.Endpoint     `validate:"required"`
	EndpointRules []entities.EndpointRule `validate:"required"`
}

type WorkspaceImportRes struct {
	WorkspaceIds     []string
	WorkspaceTierIds []string
	ApplicationIds   []string
	EndpointIds      []string
	EndpointRuleIds  []string
}

type WorkspaceExportReq struct {
	WorkspaceIds   []string `validate:"required,gt=0,dive,startswith=ws_"`
	ApplicationIds []string `validate:"required,dive,startswith=app_"`
	EndpointIds    []string `validate:"required,dive,startswith=ep_"`
}

type WorkspaceExportRes struct {
}

type WorkspaceUpdateReq struct {
	Id   string `validate:"required,startswith=ws_"`
	Name string `validate:"required"`
}

type WorkspaceUpdateRes struct {
	Doc *entities.Workspace
}

type WorkspaceGetReq struct {
	Id string `validate:"required,startswith=ws_"`
}

type WorkspaceGetRes struct {
	Workspace     *entities.Workspace
	WorkspaceTier *entities.WorkspaceTier
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metrics.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
