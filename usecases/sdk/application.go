package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type Application interface {
	Create(ctx context.Context, req *ApplicationCreateReq) (*ApplicationCreateRes, error)
	Update(ctx context.Context, req *ApplicationUpdateReq) (*ApplicationUpdateRes, error)
	Delete(ctx context.Context, req *ApplicationDeleteReq) (*ApplicationDeleteRes, error)

	List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error)
	Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error)
}

type ApplicationCreateReq struct {
	Name string `validate:"required"`
}

type ApplicationCreateRes struct {
	Doc *entities.Application
}

type ApplicationUpdateReq struct {
	Id   string `validate:"required"`
	Name string `validate:"required"`
}

type ApplicationUpdateRes struct {
	Doc *entities.Application
}

type ApplicationDeleteReq struct {
	Id string `validate:"required"`
}

type ApplicationDeleteRes struct {
	Doc *entities.Application
}

type ApplicationListReq struct {
	*structure.ListReq
}

type ApplicationListRes struct {
	*structure.ListRes[entities.Application]
}

type ApplicationGetReq struct {
	Id string `validate:"required"`
}

type ApplicationGetRes struct {
	Doc *entities.Application
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
