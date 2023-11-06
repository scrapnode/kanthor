package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationGetIn struct {
	WsId string
	Id   string
}

func (in *ApplicationGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsApp),
	)
}

type ApplicationGetOut struct {
	Doc *entities.Application
}

func (uc *application) Get(ctx context.Context, in *ApplicationGetIn) (*ApplicationGetOut, error) {
	key := CacheKeyApp(in.WsId, in.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*ApplicationGetOut, error) {
		app, err := uc.repositories.Application().Get(ctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}
		return &ApplicationGetOut{Doc: app}, nil
	})
}
