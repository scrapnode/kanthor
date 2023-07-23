package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane/repos"
)

type AppCreds interface {
	Create(ctx context.Context, req *AppCredsCreateReq) (*AppCredsCreateRes, error)
}

type AppCredsCreateReq struct {
	AppId       string                    `json:"app_id" validate:"required"`
	Role        string                    `json:"role" validate:"required"`
	Permissions []authorizator.Permission `json:"permissions" validate:"required,gt=0,dive,gt=0,dive,required,len=2"`
}

type AppCredsCreateRes struct {
	Sub   string `json:"sub"`
	Token string `json:"token"`
}

type AppCredsListReq struct {
	AppId string `json:"app_id" validate:"required"`
}

type AppCredsListRes struct {
	Data []AppCredsListEntity
}

type AppCredsListEntity struct {
	Sub         string                    `json:"sub"`
	Permissions []authorizator.Permission `json:"permissions"`
}

type appcreds struct {
	conf         *config.Config
	logger       logging.Logger
	symmetric    cryptography.Symmetric
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories
}
