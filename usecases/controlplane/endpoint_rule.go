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

type EndpointRule interface {
	List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error)
	Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error)
}

type EndpointRuleListReq struct {
	Workspace *entities.Workspace `json:"workspace" validate:"required"`
	AppId     string              `json:"app_id" validate:"required"`
	EpId      string              `json:"ep_id" validate:"required"`
	structure.ListReq
}

type EndpointRuleListRes structure.ListRes[entities.EndpointRule]

type EndpointRuleGetReq struct {
	Workspace *entities.Workspace `json:"workspace" validate:"required"`
	AppId     string              `json:"app_id" validate:"required"`
	EpId      string              `json:"ep_id" validate:"required"`
	Id        string              `json:"id" validate:"required"`
}

type EndpointRuleGetRes struct {
	EndpointRule *entities.EndpointRule `json:"endpoint"`
}

type endpointRule struct {
	conf         *config.Config
	logger       logging.Logger
	symmetric    cryptography.Symmetric
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories
}
