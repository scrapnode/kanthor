package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointDeleteIn struct {
	WsId string
	Id   string
}

func (in *EndpointDeleteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEp),
	)
}

type EndpointDeleteOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error) {
	ep, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ep, err := uc.repositories.Endpoint().Get(ctx, in.WsId, in.Id)
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
