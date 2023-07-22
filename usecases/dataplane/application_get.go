package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (usecase *application) Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error) {
	res, err := cache.Warp(usecase.cache,
		cache.Key("APPLICATION_WITH_WORKSPACE", req.Id),
		time.Hour*24,
		func() (*ApplicationGetRes, error) {
			usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_application_get"))

			app, err := usecase.repos.Application().Get(ctx, req.Id)
			if err != nil {
				return nil, err
			}
			ws, err := usecase.repos.Workspace().Get(ctx, req.Id)
			if err != nil {
				return nil, err
			}

			return &ApplicationGetRes{Application: app, Workspace: ws}, nil
		},
	)

	return res, err
}
