package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointGetReq struct {
	WsId  string
	AppId string
	Id    string
}

func (req *EndpointGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", req.Id, entities.IdNsEp),
	)
}

type EndpointGetRes struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Get(ctx context.Context, req *EndpointGetReq) (*EndpointGetRes, error) {
	key := CacheKeyEp(req.AppId, req.Id)
	// @TODO: remove hardcode time-to-live
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*EndpointGetRes, error) {
		app, err := uc.repos.Application().Get(ctx, req.WsId, req.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repos.Endpoint().Get(ctx, app, req.Id)
		if err != nil {
			return nil, err
		}
		res := &EndpointGetRes{Doc: ep}
		return res, nil
	})
}
