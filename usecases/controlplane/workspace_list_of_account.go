package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *workspace) ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error) {
	cacheKey := cache.Key("WORKSPACES_OF_ACCOUNT", req.Account.Sub)
	list, err := cache.Warp(uc.cache, cacheKey, time.Hour*24, func() (*structure.ListRes[entities.Workspace], error) {
		uc.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_list_of_accounts"))
		return uc.repos.Workspace().List(ctx, structure.WithListIds(req.WorkspaceIds))
	})
	if err != nil {
		return nil, err
	}

	res := &WorkspaceListOfAccountRes{Workspaces: list.Data}
	return res, nil
}
