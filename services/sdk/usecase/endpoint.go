package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type Endpoint interface {
	Create(ctx context.Context, req *EndpointCreateReq) (*EndpointCreateRes, error)
	Update(ctx context.Context, req *EndpointUpdateReq) (*EndpointUpdateRes, error)
	Delete(ctx context.Context, req *EndpointDeleteReq) (*EndpointDeleteRes, error)

	List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error)
	Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error)
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
