package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *endpointRule) Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error) {
	ep, err := cache.Warp(uc.cache,
		cache.Key("ENDPOINT_RULE", req.Workspace.Id, req.AppId, req.EpId, req.Id),
		time.Hour*24,
		func() (*entities.EndpointRule, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "controlplane_endpoint_rule_get"))

			return uc.repos.EndpointRule().Get(ctx, req.Workspace.Id, req.AppId, req.EpId, req.Id)
		},
	)
	if err != nil {
		return nil, err
	}

	return &EndpointRuleGetRes{EndpointRule: ep}, nil
}
