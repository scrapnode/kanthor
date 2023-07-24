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

type Endpoint interface {
	List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error)
	Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error)
}

type EndpointListReq struct {
	Workspace *entities.Workspace `json:"workspace" validate:"required"`
	AppId     string              `json:"app_id" validate:"required"`
	structure.ListReq
}

type EndpointListRes structure.ListRes[entities.Endpoint]

type EndpointGetReq struct {
	Workspace *entities.Workspace `json:"workspace" validate:"required"`
	AppId     string              `json:"app_id" validate:"required"`
	Id        string              `json:"id" validate:"required"`
}

type EndpointGetRes struct {
	Endpoint *entities.Endpoint      `json:"endpoint"`
	Rules    []entities.EndpointRule `json:"rules"`
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	symmetric    cryptography.Symmetric
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories
}
