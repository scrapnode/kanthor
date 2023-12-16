package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleUpdateIn struct {
	Ws   *entities.Workspace
	EpId string
	Id   string
	Name string
}

func (in *EndpointRuleUpdateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.PointerNotNil("ws", in.Ws),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("id", in.EpId, entities.IdNsEpr),
		validator.StringRequired("name", in.Name),
	)
}

type EndpointRuleUpdateOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Update(ctx context.Context, in *EndpointRuleUpdateIn) (*EndpointRuleUpdateOut, error) {
	epr, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ep, err := uc.repositories.Endpoint().GetOfWorkspace(txctx, in.Ws, in.EpId)
		if err != nil {
			return nil, err
		}
		epr, err := uc.repositories.EndpointRule().Get(txctx, ep, in.Id)
		if err != nil {
			return nil, err
		}

		epr.Name = in.Name
		epr.SetAT(uc.infra.Timer.Now())
		return uc.repositories.EndpointRule().Update(txctx, epr)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleUpdateOut{Doc: epr.(*entities.EndpointRule)}
	return res, nil
}
