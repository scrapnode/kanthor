package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
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
	AppId string `json:"app_id" validate:"required,startswith=app_"`
	Name  string `json:"name" validate:"required"`

	SecretKey string `json:"secret_key" validate:"omitempty,min=16,max=32"`
	Uri       string `json:"uri" validate:"required,uri"`
	Method    string `json:"method" validate:"required,oneof=POST PUT"`
}

type EndpointCreateRes struct {
	Doc *entities.Endpoint
}

type EndpointUpdateReq struct {
	AppId string `json:"app_id" validate:"required,startswith=app_"`
	Id    string `validate:"required"`
	Name  string `json:"name" validate:"required"`
}

type EndpointUpdateRes struct {
	Doc *entities.Endpoint
}

type EndpointDeleteReq struct {
	AppId string `json:"app_id" validate:"required,startswith=app_"`
	Id    string `validate:"required"`
}

type EndpointDeleteRes struct {
	Doc *entities.Endpoint
}

type EndpointListReq struct {
	AppId string `json:"app_id" validate:"required,startswith=app_"`
	*structure.ListReq
}

type EndpointListRes struct {
	*structure.ListRes[entities.Endpoint]
}

type EndpointGetReq struct {
	AppId string `json:"app_id" validate:"required,startswith=app_"`
	Id    string `validate:"required"`
}

type EndpointGetRes struct {
	Doc *entities.Endpoint
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
