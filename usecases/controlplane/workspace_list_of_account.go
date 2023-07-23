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
	// the list workspace of account is infrequently changed, so we can cache it to improve perf
	cacheKey := cache.Key("WORKSPACES_OF_ACCOUNT", req.Account.Sub)
	list, err := cache.Warp(uc.cache, cacheKey, time.Hour*24, func() (*structure.ListRes[entities.Workspace], error) {
		uc.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_list_of_account"))

		result := &structure.ListRes[entities.Workspace]{Data: []entities.Workspace{}}

		// get assigned workspaces
		if len(req.AssignedWorkspaceIds) > 0 {
			assigned, err := uc.repos.Workspace().List(ctx, structure.WithListIds(req.AssignedWorkspaceIds))
			if err != nil {
				return nil, err
			}
			result.Data = append(result.Data, assigned.Data...)
		}

		// get the workspace that account owns
		ws, err := uc.repos.Workspace().GetOwned(ctx, req.Account.Sub)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, *ws)

		return result, nil
	})
	if err != nil {
		return nil, err
	}

	res := &WorkspaceListOfAccountRes{Workspaces: list.Data}
	return res, nil
}
