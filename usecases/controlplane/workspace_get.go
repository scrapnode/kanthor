package controlplane

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (usecase *worksapce) Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error) {
	cacheKey := cache.Key("WORKSPACE", req.Id)
	ws, err := cache.Warp(usecase.cache, cacheKey, time.Hour*24, func() (*entities.Workspace, error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_workspace_get"))
		return usecase.repos.Workspace().Get(ctx, req.Id)
	})
	if err != nil {
		usecase.logger.Errorw(err.Error(), "ws_id", req.Id)
		return nil, fmt.Errorf("unable to find workspace [%s]", req.Id)
	}

	res := &WorkspaceGetRes{Workspace: ws}
	return res, nil
}
