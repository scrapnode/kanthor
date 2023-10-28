package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationDeleteReq struct {
	WsId string
	Id   string
}

func (req *ApplicationDeleteReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsApp),
	)
}

type ApplicationDeleteRes struct {
	Doc *entities.Application
}

func (uc *application) Delete(ctx context.Context, req *ApplicationDeleteReq) (*ApplicationDeleteRes, error) {
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Application().Get(txctx, req.WsId, req.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repos.Application().Delete(txctx, app); err != nil {
			return nil, err
		}
		return app, nil
	})
	if err != nil {
		return nil, err
	}

	res := &ApplicationDeleteRes{Doc: app.(*entities.Application)}
	return res, nil
}
