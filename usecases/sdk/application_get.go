package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationGetReq struct {
	Id string
}

func (req *ApplicationGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", req.Id, entities.IdNsApp),
	)
}

type ApplicationGetRes struct {
	Doc *entities.Application
}

func (uc *application) Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	key := CacheKeyApp(ws.Id, req.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*ApplicationGetRes, error) {
		app, err := uc.repos.Application().Get(ctx, ws, req.Id)
		if err != nil {
			return nil, err
		}
		res := &ApplicationGetRes{Doc: app}
		return res, nil
	})
}
