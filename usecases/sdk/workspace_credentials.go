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

type WorkspaceCredentials interface {
	Authenticate(ctx context.Context, req *WorkspaceCredentialsAuthenticateReq) (*WorkspaceCredentialsAuthenticateRes, error)
	Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error)
}

type workspaceCredentials struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
