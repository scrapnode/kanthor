package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *endpointRule) Create(ctx context.Context, req *EndpointRuleCreateReq) (*EndpointRuleCreateRes, error) {
	doc := &entities.EndpointRule{
		EndpointId:          req.EpId,
		Name:                req.Name,
		Priority:            req.Priority,
		Exclusionary:        req.Exclusionary,
		ConditionSource:     req.ConditionSource,
		ConditionExpression: req.ConditionExpression,
	}
	doc.GenId()
	doc.SetAT(uc.timer.Now())

	ep, err := uc.repos.EndpointRule().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleCreateRes{Doc: ep}
	return res, nil
}
