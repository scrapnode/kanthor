package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repos"
)

type Trigger interface {
	Plan(ctx context.Context, req *TriggerPlanReq) (*TriggerPlanRes, error)
	Exec(ctx context.Context, req *TriggerExecReq) (*TriggerExecRes, error)
}

type trigger struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
