package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointGetReq struct {
	AppId string
	Id    string
}

func (req *EndpointGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, "app_"),
		validator.StringStartsWith("id", req.Id, "ep_"),
	)
}

type EndpointGetRes struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	key := CacheKeyEp(req.AppId, req.Id)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*EndpointGetRes, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)
		app, err := uc.repos.Application().Get(ctx, ws, req.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repos.Endpoint().Get(ctx, app, req.Id)
		if err != nil {
			return nil, err
		}
		res := &EndpointGetRes{Doc: ep}
		return res, nil
	})
}
