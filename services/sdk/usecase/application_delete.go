package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationDeleteIn struct {
	WsId string
	Id   string
}

func (in *ApplicationDeleteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsApp),
	)
}

type ApplicationDeleteOut struct {
	Doc *entities.Application
}

func (uc *application) Delete(ctx context.Context, in *ApplicationDeleteIn) (*ApplicationDeleteOut, error) {
	app, err := uc.repositories.Database().Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repositories.Database().Application().Get(txctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repositories.Database().Application().Delete(txctx, app); err != nil {
			return nil, err
		}
		return app, nil
	})
	if err != nil {
		return nil, err
	}

	return &ApplicationDeleteOut{Doc: app.(*entities.Application)}, nil
}
