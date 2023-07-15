package controlplane

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (usecase *worksapce) ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error) {
	cacheKey := cache.Key("WORKSPACES_OF_ACCOUNT", req.Account.Sub)
	list, err := cache.Warp(usecase.cache, cacheKey, time.Hour*24, func() (*structure.ListRes[entities.Workspace], error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_list_of_accounts"))
		return usecase.repos.Workspace().List(
			ctx,
			structure.WithListCursor(req.Cursor),
			structure.WithListSearch(req.Search),
			structure.WithListLimit(req.Limit),
		)
	})
	if err != nil {
		usecase.logger.Errorw(err.Error(), "account_sub", req.Account.Sub)
		return nil, fmt.Errorf("unable to list workspace of account [%v]", req.Account.Sub)
	}

	res := &WorkspaceListOfAccountRes{Workspaces: list.Data}
	return res, nil
}
