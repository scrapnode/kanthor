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

type EndpointRule interface {
	Create(ctx context.Context, req *EndpointRuleCreateReq) (*EndpointRuleCreateRes, error)
	Update(ctx context.Context, req *EndpointRuleUpdateReq) (*EndpointRuleUpdateRes, error)
	Delete(ctx context.Context, req *EndpointRuleDeleteReq) (*EndpointRuleDeleteRes, error)

	List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error)
	Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error)
}

type endpointRule struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
