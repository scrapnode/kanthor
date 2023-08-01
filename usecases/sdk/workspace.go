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

type Workspace interface {
	Authenticate(ctx context.Context, req *WorkspaceAuthenticateReq) (*WorkspaceAuthenticateRes, error)
}

type WorkspaceAuthenticateReq struct {
	User string `validate:"required"`
	Hash string `validate:"required"`
}

type WorkspaceAuthenticateRes struct {
	Error                error
	Workspace            *entities.Workspace
	WorkspaceCredentials *entities.WorkspaceCredentials
	WorkspaceTier        *entities.WorkspaceTier
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
