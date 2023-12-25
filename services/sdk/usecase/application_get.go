package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
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
	app, err := uc.repositories.Application().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}
	return &ApplicationGetOut{Doc: app}, nil
}
