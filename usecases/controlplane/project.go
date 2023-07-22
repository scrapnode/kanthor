package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
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
	Account       *authenticator.Account `json:"account" validate:"required"`
	WorkspaceName string                 `json:"workspace_name" validate:"required"`
	WorkspaceTier string                 `json:"workspace_tier" validate:"required"`
}

type ProjectSetupDefaultRes struct {
	WorkspaceId   string `json:"workspace_id"`
	WorkspaceTier string `json:"workspace_tier"`
}

type ProjectSetupDemoReq struct {
	Account     *authenticator.Account `json:"account" validate:"required"`
	WorkspaceId string                 `json:"workspace_id" validate:"required"`

	Applications  []entities.Application  `json:"applications" validate:"required"`
	Endpoints     []entities.Endpoint     `json:"endpoints" validate:"required"`
	EndpointRules []entities.EndpointRule `json:"endpoint_rules" validate:"required"`
}

type ProjectSetupDemoRes struct {
	WorkspaceId     string   `json:"workspace_id"`
	WorkspaceTier   string   `json:"workspace_tier"`
	ApplicationIds  []string `json:"application_ids"`
	EndpointIds     []string `json:"endpoint_ids"`
	EndpointRuleIds []string `json:"endpoint_rule_ids"`
}

type project struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
