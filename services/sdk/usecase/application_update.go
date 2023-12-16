package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationUpdateIn struct {
	WsId string
	Id   string
	Name string
}

func (in *ApplicationUpdateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsApp),
		validator.StringRequired("name", in.Name),
	)
}

type ApplicationUpdateOut struct {
	Doc *entities.Application
}

func (uc *application) Update(ctx context.Context, in *ApplicationUpdateIn) (*ApplicationUpdateOut, error) {
	app, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repositories.Application().Get(txctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}

		app.Name = in.Name
		app.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Application().Update(txctx, app)
	})
	if err != nil {
		return nil, err
	}

	return &ApplicationUpdateOut{Doc: app.(*entities.Application)}, nil
}
