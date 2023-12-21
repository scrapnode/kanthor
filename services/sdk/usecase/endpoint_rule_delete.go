package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleDeleteIn struct {
	WsId string
	Id   string
}

func (in *EndpointRuleDeleteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEpr),
	)
}

type EndpointRuleDeleteOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Delete(ctx context.Context, in *EndpointRuleDeleteIn) (*EndpointRuleDeleteOut, error) {
	epr, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		epr, err := uc.repositories.EndpointRule().Get(ctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repositories.EndpointRule().Delete(txctx, epr); err != nil {
			return nil, err
		}
		return epr, nil
	})
	if err != nil {
		return nil, err
	}

	return &EndpointRuleDeleteOut{Doc: epr.(*entities.EndpointRule)}, nil
}
