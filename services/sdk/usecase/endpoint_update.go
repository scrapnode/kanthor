package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointUpdateIn struct {
	WsId string
	Id   string
	Name string
}

func (in *EndpointUpdateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEp),
		validator.StringRequired("name", in.Name),
	)
}

type EndpointUpdateOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error) {
	ep, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ep, err := uc.repositories.Endpoint().Get(ctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}

		ep.Name = in.Name
		ep.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Endpoint().Update(txctx, ep)
	})
	if err != nil {
		return nil, err
	}

	return &EndpointUpdateOut{Doc: ep.(*entities.Endpoint)}, nil
}
