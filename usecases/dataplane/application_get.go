package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *application) Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error) {
	res, err := cache.Warp(uc.cache,
		cache.Key("APPLICATION_WITH_WORKSPACE", req.Id),
		time.Hour*24,
		func() (*ApplicationGetRes, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_application_get"))

			app, err := uc.repos.Application().Get(ctx, req.Id)
			if err != nil {
				return nil, err
			}
			ws, err := uc.repos.Workspace().Get(ctx, app.WorkspaceId)
			if err != nil {
				return nil, err
			}

			return &ApplicationGetRes{Application: app, Workspace: ws}, nil
		},
	)

	return res, err
}
