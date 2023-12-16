package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointGetIn struct {
	WsId  string
	AppId string
	Id    string
}

func (in *EndpointGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", in.Id, entities.IdNsEp),
	)
}

type EndpointGetOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error) {
	key := CacheKeyEp(in.AppId, in.Id)
	// @TODO: remove hardcode time-to-live
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*EndpointGetOut, error) {
		app, err := uc.repositories.Application().Get(ctx, in.WsId, in.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repositories.Endpoint().Get(ctx, app, in.Id)
		if err != nil {
			return nil, err
		}
		return &EndpointGetOut{Doc: ep}, nil
	})
}
