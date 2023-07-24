package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *endpoint) Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error) {
	ep, err := cache.Warp(uc.cache,
		cache.Key("ENDPOINT", req.Workspace.Id, req.AppId, req.Id),
		time.Hour*24,
		func() (*entities.Endpoint, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "controlplane_endpoint_get"))

			return uc.repos.Endpoint().Get(ctx, req.Workspace.Id, req.AppId, req.Id)
		},
	)
	if err != nil {
		return nil, err
	}

	return &EndpointGetRes{Endpoint: ep}, nil
}
