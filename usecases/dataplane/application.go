package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane/repos"
)

type Application interface {
	Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error)
}

type ApplicationGetReq struct {
	Id string `json:"id" validate:"required"`
}

type ApplicationGetRes struct {
	Application *entities.Application `json:"application"`
	Workspace   *entities.Workspace   `json:"workspace"`
}

type application struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories
}
