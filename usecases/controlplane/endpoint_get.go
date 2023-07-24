package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *endpoint) Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error) {
	res, err := cache.Warp(uc.cache,
		cache.Key("ENDPOINT", req.Workspace.Id, req.AppId, req.Id),
		time.Hour*24,
		func() (*EndpointGetRes, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "controlplane_endpoint_get"))

			ep, err := uc.repos.Endpoint().Get(ctx, req.Workspace.Id, req.AppId, req.Id)
			if err != nil {
				return nil, err
			}

			// don't put limit here because we want to get all rules of this endpoint
			rules, err := uc.repos.EndpointRule().List(ctx, req.Workspace.Id, req.AppId, ep.Id)
			if err != nil {
				return nil, err
			}

			return &EndpointGetRes{Endpoint: ep, Rules: rules.Data}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
