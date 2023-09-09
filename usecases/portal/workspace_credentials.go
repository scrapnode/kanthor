package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
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

type WorkspaceCredentialsGenerateReq struct {
	WorkspaceId string `validate:"required,startswith=ws_"`
	Name        string `validate:"required"`
	ExpiredAt   int64
}

type WorkspaceCredentialsGenerateRes struct {
	Credentials *entities.WorkspaceCredentials
	Password    string
}

type WorkspaceCredentialsUpdateReq struct {
	WorkspaceId string `validate:"required,startswith=ws_"`
	Id          string `validate:"required,startswith=wsc_"`
	Name        string `validate:"required"`
}

type WorkspaceCredentialsUpdateRes struct {
	Doc *entities.WorkspaceCredentials
}

type WorkspaceCredentialsExpireReq struct {
	WorkspaceId string `validate:"required,startswith=ws_"`
	Id          string `validate:"required,startswith=wsc_"`
	Duration    int64  `validate:"gte=0"`
}

type WorkspaceCredentialsExpireRes struct {
	Id        string
	ExpiredAt int64
}

type WorkspaceCredentialsListReq struct {
	*structure.ListReq
}

type WorkspaceCredentialsListRes struct {
	*structure.ListRes[entities.WorkspaceCredentials]
}

type WorkspaceCredentialsGetReq struct {
	WorkspaceId string `validate:"required,startswith=ws_"`
	Id          string `validate:"required,startswith=wsc_"`
}

type WorkspaceCredentialsGetRes struct {
	Doc *entities.WorkspaceCredentials
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
