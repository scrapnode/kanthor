package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleCreateIn struct {
	Ws   *entities.Workspace
	EpId string
	Name string

	Priority            int32
	Exclusionary        bool
	ConditionSource     string
	ConditionExpression string
}

func (in *EndpointRuleCreateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.PointerNotNil("ws", in.Ws),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.StringRequired("name", in.Name),
		validator.NumberGreaterThan("priority", in.Priority, 0),
		validator.StringRequired("condition_source", in.ConditionSource),
		validator.StringRequired("condition_expression", in.ConditionExpression),
	)
}

type EndpointRuleCreateOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Create(ctx context.Context, in *EndpointRuleCreateIn) (*EndpointRuleCreateOut, error) {
	ep, err := uc.repositories.Endpoint().GetOfWorkspace(ctx, in.Ws, in.EpId)
	if err != nil {
		return nil, err
	}

	doc := &entities.EndpointRule{
		EpId:                ep.Id,
		Name:                in.Name,
		Priority:            in.Priority,
		Exclusionary:        in.Exclusionary,
		ConditionSource:     in.ConditionSource,
		ConditionExpression: in.ConditionExpression,
	}
	doc.Id = suid.New(entities.IdNsEpr)
	doc.SetAT(uc.infra.Timer.Now())

	epr, err := uc.repositories.EndpointRule().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &EndpointRuleCreateOut{Doc: epr}, nil
}
