package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"time"
)

func (uc *endpointRule) Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	key := CacheKeyEpr(ws.Id, req.AppId, req.EpId, req.Id)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*EndpointRuleGetRes, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

		app, err := uc.repos.EndpointRule().Get(ctx, ws.Id, req.AppId, req.EpId, req.Id)
		if err != nil {
			return nil, err
		}
		res := &EndpointRuleGetRes{Doc: app}
		return res, nil
	})
}
