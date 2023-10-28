package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointUpdateReq struct {
	WsId  string
	AppId string
	Id    string
	Name  string
}

func (req *EndpointUpdateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", req.Id, entities.IdNsEp),
		validator.StringRequired("name", req.Name),
	)
}

type EndpointUpdateRes struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Update(ctx context.Context, req *EndpointUpdateReq) (*EndpointUpdateRes, error) {
	ep, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repositories.Application().Get(ctx, req.WsId, req.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repositories.Endpoint().Get(txctx, app, req.Id)
		if err != nil {
			return nil, err
		}

		ep.Name = req.Name
		ep.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Endpoint().Update(txctx, ep)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointUpdateRes{Doc: ep.(*entities.Endpoint)}
	return res, nil
}
