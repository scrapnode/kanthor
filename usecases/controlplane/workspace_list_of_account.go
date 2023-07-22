package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (usecase *workspace) ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error) {
	wsIds, err := usecase.authorizator.Tenants(req.Account.Sub)
	if err != nil {
		return nil, err
	}

	cacheKey := cache.Key("WORKSPACES_OF_ACCOUNT", req.Account.Sub)
	list, err := cache.Warp(usecase.cache, cacheKey, time.Hour*24, func() (*structure.ListRes[entities.Workspace], error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_list_of_accounts"))
		return usecase.repos.Workspace().List(ctx, structure.WithListIds(wsIds))
	})
	if err != nil {
		return nil, err
	}

	res := &WorkspaceListOfAccountRes{Workspaces: list.Data}
	return res, nil
}
