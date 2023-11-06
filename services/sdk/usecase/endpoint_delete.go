package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointDeleteIn struct {
	WsId  string
	AppId string
	Id    string
}

func (in *EndpointDeleteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", in.Id, entities.IdNsEp),
	)
}

type EndpointDeleteOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error) {
	ep, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repositories.Application().Get(ctx, in.WsId, in.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repositories.Endpoint().Get(txctx, app, in.Id)
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

	return &EndpointDeleteOut{Doc: ep.(*entities.Endpoint)}, nil
}
