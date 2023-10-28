package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type Application interface {
	Create(ctx context.Context, req *ApplicationCreateReq) (*ApplicationCreateRes, error)
	Update(ctx context.Context, req *ApplicationUpdateReq) (*ApplicationUpdateRes, error)
	Delete(ctx context.Context, req *ApplicationDeleteReq) (*ApplicationDeleteRes, error)

	List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error)
	Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error)
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
