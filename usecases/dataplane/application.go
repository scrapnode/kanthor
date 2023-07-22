package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/crypto"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane/repos"
)

type Application interface {
	Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error)
	GenToken(ctx context.Context, req *ApplicationGenTokenReq) (*ApplicationGenTokenRes, error)
}

type ApplicationGetReq struct {
	Id string `json:"id" validate:"required"`
}

type ApplicationGetRes struct {
	Application *entities.Application `json:"application"`
	Workspace   *entities.Workspace   `json:"workspace"`
}

type ApplicationGenTokenReq struct {
	Id          string     `json:"id" validate:"required"`
	Role        string     `json:"role" validate:"required"`
	Permissions [][]string `json:"permissions" yaml:"permissions" validate:"required,gt=0,dive,gt=0,dive,required,len=2"`
}

type ApplicationGenTokenRes struct {
	Sub   string `json:"sub"`
	Token string `json:"token"`
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	aes          *crypto.AES
	repos        repos.Repositories
}
