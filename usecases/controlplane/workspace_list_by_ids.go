package controlplane

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/utils"
	"sort"
	"time"
)

func (usecase *worksapce) ListByIds(ctx context.Context, req *WorkspaceListByIdsReq) (*WorkspaceListByIdsRes, error) {
	if len(req.Ids) == 0 {
		return &WorkspaceListByIdsRes{Workspaces: []entities.Workspace{}}, nil
	}

	// workspaces of user does not change frequently, so we can sort the list of id and use cache here
	sort.Strings(req.Ids)
	cacheKey := cache.Key("WORKSPACE_BY_IDS", utils.Key(req.Ids...))
	list, err := cache.Warp(usecase.cache, cacheKey, time.Hour*24, func() (*structure.ListRes[entities.Workspace], error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_list_by_ids"))
		return usecase.repos.Workspace().ListByIds(ctx, req.Ids)
	})
	if err != nil {
		usecase.logger.Errorw(err.Error(), "ws_ids", req.Ids)
		return nil, fmt.Errorf("unable to list workspace by ids [%v]", req.Ids)
	}

	res := &WorkspaceListByIdsRes{Workspaces: list.Data}
	return res, nil
}
