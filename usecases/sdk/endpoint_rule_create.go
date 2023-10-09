package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleCreateReq struct {
	EpId string
	Name string

	Priority            int32
	Exclusionary        bool
	ConditionSource     string
	ConditionExpression string
}

func (req *EndpointRuleCreateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", req.EpId, entities.IdNsEp),
		validator.StringRequired("name", req.Name),
		validator.NumberGreaterThan("priority", req.Priority, 0),
		validator.StringRequired("condition_source", req.ConditionSource),
		validator.StringRequired("condition_expression", req.ConditionExpression),
	)
}

type EndpointRuleCreateRes struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Create(ctx context.Context, req *EndpointRuleCreateReq) (*EndpointRuleCreateRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	ep, err := uc.repos.Endpoint().GetOfWorkspace(ctx, ws, req.EpId)
	if err != nil {
		return nil, err
	}

	doc := &entities.EndpointRule{
		EpId:                ep.Id,
		Name:                req.Name,
		Priority:            req.Priority,
		Exclusionary:        req.Exclusionary,
		ConditionSource:     req.ConditionSource,
		ConditionExpression: req.ConditionExpression,
	}
	doc.GenId()
	doc.SetAT(uc.timer.Now())

	epr, err := uc.repos.EndpointRule().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleCreateRes{Doc: epr}
	return res, nil
}
