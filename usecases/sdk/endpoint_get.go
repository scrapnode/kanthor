package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"time"
)

func (uc *endpoint) Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	key := CacheKeyEp(ws.Id, req.AppId, req.Id)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*EndpointGetRes, error) {
		uc.metrics.Count("cache_miss_total", 1)

		app, err := uc.repos.Endpoint().Get(ctx, ws.Id, req.AppId, req.Id)
		if err != nil {
			return nil, err
		}
		res := &EndpointGetRes{Doc: app}
		return res, nil
	})
}
