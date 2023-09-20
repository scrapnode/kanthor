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

type Application interface {
	Create(ctx context.Context, req *ApplicationCreateReq) (*ApplicationCreateRes, error)
	Update(ctx context.Context, req *ApplicationUpdateReq) (*ApplicationUpdateRes, error)
	Delete(ctx context.Context, req *ApplicationDeleteReq) (*ApplicationDeleteRes, error)

	List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error)
	Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error)
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
