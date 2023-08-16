package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type WorkspaceCredentials interface {
	Authenticate(ctx context.Context, req *WorkspaceCredentialsAuthenticateReq) (*WorkspaceCredentialsAuthenticateRes, error)
	Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error)
}

type WorkspaceCredentialsAuthenticateReq struct {
	User string `validate:"required,startswith=wsc_"`
	Hash string `validate:"required"`
}

type WorkspaceCredentialsAuthenticateRes struct {
	Workspace            *entities.Workspace
	WorkspaceCredentials *entities.WorkspaceCredentials
	WorkspaceTier        *entities.WorkspaceTier
}

type WorkspaceCredentialsExpireReq struct {
	User      string `validate:"required,startswith=wsc_"`
	ExpiredAt int64  `validate:"required,gt=0"`
}

type WorkspaceCredentialsExpireRes struct {
	Ok bool
}

type workspaceCredentials struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
