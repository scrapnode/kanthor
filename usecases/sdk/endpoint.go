package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type Endpoint interface {
	Create(ctx context.Context, req *EndpointCreateReq) (*EndpointCreateRes, error)
	Update(ctx context.Context, req *EndpointUpdateReq) (*EndpointUpdateRes, error)
	Delete(ctx context.Context, req *EndpointDeleteReq) (*EndpointDeleteRes, error)

	List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error)
	Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error)
}

type EndpointCreateReq struct {
	AppId string `validate:"required,startswith=app_"`
	Name  string `validate:"required"`

	SecretKey string `validate:"omitempty,min=16,max=32"`
	Uri       string `validate:"required,uri"`
	Method    string `validate:"required,oneof=POST PUT"`
}

type EndpointCreateRes struct {
	Doc *entities.Endpoint
}

type EndpointUpdateReq struct {
	AppId string `validate:"required,startswith=app_"`
	Id    string `validate:"required,startswith=ep_"`
	Name  string `validate:"required"`
}

type EndpointUpdateRes struct {
	Doc *entities.Endpoint
}

type EndpointDeleteReq struct {
	AppId string `validate:"required,startswith=app_"`
	Id    string `validate:"required,startswith=ep_"`
}

type EndpointDeleteRes struct {
	Doc *entities.Endpoint
}

type EndpointListReq struct {
	AppId string `validate:"required,startswith=app_"`
	*structure.ListReq
}

type EndpointListRes struct {
	*structure.ListRes[entities.Endpoint]
}

type EndpointGetReq struct {
	AppId string `validate:"required,startswith=app_"`
	Id    string `validate:"required,startswith=ep_"`
}

type EndpointGetRes struct {
	Doc *entities.Endpoint
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metrics.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
