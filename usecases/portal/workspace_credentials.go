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

type WorkspaceCredentials interface {
	Generate(ctx context.Context, req *WorkspaceCredentialsReq) (*WorkspaceCredentialsRes, error)
}

type WorkspaceCredentialsReq struct {
	WorkspaceId string `validate:"required"`
	Count       int    `validate:"required,gt=0,lt=10"`
}

type WorkspaceCredentialsRes struct {
	Credentials []entities.WorkspaceCredentials
}

type workspaceCredentials struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
