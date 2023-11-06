package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleDeleteIn struct {
	EpId string
	Id   string
}

func (in *EndpointRuleDeleteIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("id", in.EpId, entities.IdNsEpr),
	)
}

type EndpointRuleDeleteOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Delete(ctx context.Context, in *EndpointRuleDeleteIn) (*EndpointRuleDeleteOut, error) {
	ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

	epr, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ep, err := uc.repositories.Endpoint().GetOfWorkspace(txctx, ws, in.EpId)
		if err != nil {
			return nil, err
		}

		epr, err := uc.repositories.EndpointRule().Get(txctx, ep, in.Id)
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
