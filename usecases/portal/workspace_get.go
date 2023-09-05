package portal

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *workspace) Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error) {
	key := utils.Key("portal", req.Id)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*WorkspaceGetRes, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

		ws, err := uc.repos.Workspace().Get(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		tier, err := uc.repos.WorkspaceTier().Get(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		res := &WorkspaceGetRes{Workspace: ws, WorkspaceTier: tier}
		return res, nil
	})
}
