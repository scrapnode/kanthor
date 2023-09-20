package portal

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type WorkspaceCredentials interface {
	Generate(ctx context.Context, req *WorkspaceCredentialsGenerateReq) (*WorkspaceCredentialsGenerateRes, error)
	Update(ctx context.Context, req *WorkspaceCredentialsUpdateReq) (*WorkspaceCredentialsUpdateRes, error)
	Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error)
	List(ctx context.Context, req *WorkspaceCredentialsListReq) (*WorkspaceCredentialsListRes, error)
	Get(ctx context.Context, req *WorkspaceCredentialsGetReq) (*WorkspaceCredentialsGetRes, error)
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
