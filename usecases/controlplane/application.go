package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

type Application interface {
	List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error)
	Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error)
}

type ApplicationListReq struct {
	Workspace *entities.Workspace `json:"workspace" validate:"required"`
	structure.ListReq
}

type ApplicationListRes structure.ListRes[entities.Application]

type ApplicationGetReq struct {
	Workspace *entities.Workspace `json:"workspace" validate:"required"`
	Id        string              `json:"id" validate:"required"`
}

type ApplicationGetRes struct {
	Application *entities.Application `json:"application"`
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	symmetric    cryptography.Symmetric
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories
}
