package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type WorkspaceCredentials interface {
	Generate(ctx context.Context, req *WorkspaceCredentialsGenerateReq) (*WorkspaceCredentialsGenerateRes, error)
}

type WorkspaceCredentialsGenerateReq struct {
	WorkspaceId string `validate:"required"`
	Count       int    `validate:"required,gt=0,lt=10"`
}

type WorkspaceCredentialsGenerateRes struct {
	Credentials []entities.WorkspaceCredentials
	Passwords   map[string]string
}

type workspaceCredentials struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
