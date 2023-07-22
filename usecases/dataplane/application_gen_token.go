package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *application) GenToken(ctx context.Context, req *ApplicationGenTokenReq) (*ApplicationGenTokenRes, error) {
	app, err := cache.Warp(uc.cache,
		cache.Key("APPLICATION", req.Id),
		time.Hour*24,
		func() (*entities.Application, error) {
			uc.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_application_gen_token_get"))

			return uc.repos.Application().Get(ctx, req.Id)
		},
	)
	if err != nil {
		return nil, err
	}

	// setup permission of a role in request app
	if err := uc.authorizator.SetupPermissions(req.Role, app.Id, req.Permissions); err != nil {
		return nil, err
	}

	// generate the account
	sub := utils.ID("appsub")
	token, err := uc.aes.EncryptString(sub)
	if err != nil {
		return nil, err
	}

	if err := uc.authorizator.GrantAccess(sub, req.Role, app.Id); err != nil {
		return nil, err
	}

	res := &ApplicationGenTokenRes{Sub: sub, Token: token}
	return res, nil
}
