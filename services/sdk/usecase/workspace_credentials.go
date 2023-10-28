package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repos"
)

type WorkspaceCredentials interface {
	Authenticate(ctx context.Context, req *WorkspaceCredentialsAuthenticateReq) (*WorkspaceCredentialsAuthenticateRes, error)
}

type workspaceCredentials struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
