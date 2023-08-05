package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *endpoint) Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	key := utils.Key(ws.Id, req.AppId, req.Id)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*EndpointGetRes, error) {
		app, err := uc.repos.Endpoint().Get(ctx, ws.Id, req.AppId, req.Id)
		if err != nil {
			return nil, err
		}
		res := &EndpointGetRes{Doc: app}
		return res, nil
	})
}
