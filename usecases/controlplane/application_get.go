package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *application) Get(ctx context.Context, req *ApplicationGetReq) (*ApplicationGetRes, error) {
	app, err := cache.Warp(uc.cache,
		cache.Key("APPLICATION", req.Workspace.Id, req.Id),
		time.Hour*24,
		func() (*entities.Application, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "controlplane_application_get"))

			return uc.repos.Application().Get(ctx, req.Workspace.Id, req.Id)
		},
	)
	if err != nil {
		return nil, err
	}

	res := &ApplicationGetRes{Application: app}
	return res, nil
}
