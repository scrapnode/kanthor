package controlplane

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (usecase *worksapce) GetByAccount(ctx context.Context, req *WorkspaceGetByAccountReq) (*WorkspaceGetByAccountRes, error) {
	cacheKey := cache.Key("WORKSPACE_BY_ACCOUNT", req.WorkspaceId, req.AccountSub)
	ws, err := cache.Warp(usecase.cache, cacheKey, time.Hour*24, func() (*entities.Workspace, error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_list_by_ids"))
		return usecase.repos.Workspace().GetByAccountSub(ctx, req.WorkspaceId, req.AccountSub)
	})
	if err != nil {
		usecase.logger.Errorw(err.Error(), "account_sub", req.AccountSub)
		return nil, fmt.Errorf("unable to list workspace of account [%v]", req.AccountSub)
	}

	res := &WorkspaceGetByAccountRes{Workspace: ws}
	return res, nil
}
