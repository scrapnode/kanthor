package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/data/demo"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

type Project interface {
	SetupDefault(ctx context.Context, req *ProjectSetupDefaultReq) (*ProjectSetupDefaultRes, error)
	SetupDemo(ctx context.Context, req *ProjectSetupDemoReq) (*ProjectSetupDemoRes, error)
}

type ProjectSetupDefaultReq struct {
	Account       *authenticator.Account
	WorkspaceName string
	WorkspaceTier string
}

type ProjectSetupDefaultRes struct {
	WorkspaceId   string
	WorkspaceTier string
}

type ProjectSetupDemoReq struct {
	Account     *authenticator.Account
	WorkspaceId string
	Entities    *demo.ProjectEntities
}

type ProjectSetupDemoRes struct {
	ApplicationIds  []string
	EndpointIds     []string
	EndpointRuleIds []string
}

type project struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
