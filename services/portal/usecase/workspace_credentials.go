package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
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
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
