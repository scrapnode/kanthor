package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleUpdateIn struct {
	WsId string
	Id   string
	Name string

	Priority            int32
	Exclusionary        bool
	ConditionSource     string
	ConditionExpression string
}

func (in *EndpointRuleUpdateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEpr),
		validator.StringRequired("name", in.Name),
		validator.NumberGreaterThan("priority", in.Priority, 0),
		validator.StringRequired("condition_source", in.ConditionSource),
		validator.StringRequired("condition_expression", in.ConditionExpression),
	)
}

type EndpointRuleUpdateOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Update(ctx context.Context, in *EndpointRuleUpdateIn) (*EndpointRuleUpdateOut, error) {
	epr, err := uc.repositories.Database().Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		epr, err := uc.repositories.Database().EndpointRule().Get(ctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}

		epr.Name = in.Name
		epr.Priority = in.Priority
		epr.Exclusionary = in.Exclusionary
		epr.ConditionSource = in.ConditionSource
		epr.ConditionExpression = in.ConditionExpression
		epr.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Database().EndpointRule().Update(txctx, epr)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleUpdateOut{Doc: epr.(*entities.EndpointRule)}
	return res, nil
}
