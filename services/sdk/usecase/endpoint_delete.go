package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointDeleteReq struct {
	WsId  string
	AppId string
	Id    string
}

func (req *EndpointDeleteReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", req.Id, entities.IdNsEp),
	)
}

type EndpointDeleteRes struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Delete(ctx context.Context, req *EndpointDeleteReq) (*EndpointDeleteRes, error) {
	ep, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repositories.Application().Get(ctx, req.WsId, req.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repositories.Endpoint().Get(txctx, app, req.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repositories.Endpoint().Delete(txctx, ep); err != nil {
			return nil, err
		}
		return ep, nil
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointDeleteRes{Doc: ep.(*entities.Endpoint)}
	return res, nil
}
