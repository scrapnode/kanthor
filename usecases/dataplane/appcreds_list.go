package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func (uc *appcreds) List(ctx context.Context, req *AppCredsListReq) (*AppCredsListRes, error) {
	res, err := cache.Warp(uc.cache,
		cache.Key("APPLICATION_CREDENTIALS", req.AppId),
		time.Hour*24,
		func() (*AppCredsListRes, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_appcreds_list"))

			res := &AppCredsListRes{Data: []AppCredsListEntity{}}

			subs, err := uc.authorizator.UsersOfTenant(req.AppId)
			if err != nil {
				return nil, err
			}

			for _, sub := range subs {
				entity := AppCredsListEntity{Sub: sub}
				permissions, err := uc.authorizator.UserPermissionsInTenant(req.AppId, sub)
				if err != nil {
					return nil, err
				}
				entity.Permissions = permissions

				res.Data = append(res.Data, entity)
			}

			return res, nil
		},
	)

	return res, err
}
