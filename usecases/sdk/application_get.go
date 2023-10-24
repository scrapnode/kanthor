package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationGetReq struct {
	WorkspaceId string
	Id          string
}

func (req *ApplicationGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WorkspaceId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsApp),
	)
}

type ApplicationGetRes struct {
	Doc *entities.Application
}

func (uc *application) Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error) {
	key := CacheKeyApp(req.WorkspaceId, req.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*ApplicationGetRes, error) {
		app, err := uc.repos.Application().Get(ctx, req.WorkspaceId, req.Id)
		if err != nil {
			return nil, err
		}
		res := &ApplicationGetRes{Doc: app}
		return res, nil
	})
}
