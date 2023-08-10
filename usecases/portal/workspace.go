package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Workspace interface {
	Import(ctx context.Context, req *WorkspaceImportReq) (*WorkspaceImportRes, error)
	Export(ctx context.Context, req *WorkspaceExportReq) (*WorkspaceExportRes, error)

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
	WorkspaceIds   []string `validate:"required,gt=0"`
	ApplicationIds []string `validate:"required"`
	EndpointIds    []string `validate:"required"`
}

type WorkspaceExportRes struct {
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
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
