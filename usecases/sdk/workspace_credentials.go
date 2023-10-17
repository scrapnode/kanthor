package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type WorkspaceCredentials interface {
	Authenticate(ctx context.Context, req *WorkspaceCredentialsAuthenticateReq) (*WorkspaceCredentialsAuthenticateRes, error)
	Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error)
}

type workspaceCredentials struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
