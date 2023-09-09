package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Workspace interface {
	Setup(ctx context.Context, req *WorkspaceSetupReq) (*WorkspaceSetupRes, error)

	Update(ctx context.Context, req *WorkspaceUpdateReq) (*WorkspaceUpdateRes, error)
	Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error)
}

type WorkspaceSetupReq struct {
	Workspace     *entities.Workspace     `validate:"required"`
	Applications  []entities.Application  `validate:"required"`
	Endpoints     []entities.Endpoint     `validate:"required"`
	EndpointRules []entities.EndpointRule `validate:"required"`
}

type WorkspaceSetupRes struct {
	ApplicationIds  []string
	EndpointIds     []string
	EndpointRuleIds []string
	Status          map[string]bool
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
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
