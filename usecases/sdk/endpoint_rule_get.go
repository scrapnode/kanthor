package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *endpointRule) Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	key := utils.Key(ws.Id, req.AppId, req.EpId, req.Id)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*EndpointRuleGetRes, error) {
		app, err := uc.repos.EndpointRule().Get(ctx, ws.Id, req.AppId, req.EpId, req.Id)
		if err != nil {
			return nil, err
		}
		res := &EndpointRuleGetRes{Doc: app}
		return res, nil
	})
}
