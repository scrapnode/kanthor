package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
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

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
