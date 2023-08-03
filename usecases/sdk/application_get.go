package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *application) Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	key := utils.Key(ws.Id, req.Id)
	return cache.Warp(uc.cache, key, time.Hour*24, func() (*ApplicationGetRes, error) {
		app, err := uc.repos.Application().Get(ctx, ws.Id, req.Id)
		if err != nil {
			return nil, err
		}
		res := &ApplicationGetRes{Doc: app}
		return res, nil
	})
}
