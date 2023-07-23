package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *appcreds) Create(ctx context.Context, req *AppCredsCreateReq) (*AppCredsCreateRes, error) {
	app, err := cache.Warp(uc.cache,
		cache.Key("APPLICATION", req.AppId),
		time.Hour*24,
		func() (*entities.Application, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_appcreds_application_get"))

			return uc.repos.Application().Get(ctx, req.AppId)
		},
	)
	if err != nil {
		return nil, err
	}

	// setup permission of a role in request app
	if err := uc.authorizator.GrantPermissionsToRole(app.Id, req.Role, req.Permissions); err != nil {
		return nil, err
	}

	// generate the account
	sub := utils.ID("appcreds")
	token, err := uc.symmetric.StringEncrypt(sub)
	if err != nil {
		return nil, err
	}

	if err := uc.authorizator.GrantRoleToSub(app.Id, sub, req.Role); err != nil {
		return nil, err
	}

	res := &AppCredsCreateRes{Sub: sub, Token: token}
	return res, nil
}
