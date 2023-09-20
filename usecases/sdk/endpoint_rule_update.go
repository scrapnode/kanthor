package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleUpdateReq struct {
	EpId string
	Id   string
	Name string
}

func (req *EndpointRuleUpdateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", req.EpId, "ep_"),
		validator.StringStartsWith("id", req.EpId, "epr_"),
		validator.StringRequired("name", req.Name),
	)
}

type EndpointRuleUpdateRes struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Update(ctx context.Context, req *EndpointRuleUpdateReq) (*EndpointRuleUpdateRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	epr, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ep, err := uc.repos.Endpoint().GetOfWorkspace(txctx, ws, req.EpId)
		if err != nil {
			return nil, err
		}
		epr, err := uc.repos.EndpointRule().Get(txctx, ep, req.Id)
		if err != nil {
			return nil, err
		}

		epr.Name = req.Name
		epr.SetAT(uc.timer.Now())
		return uc.repos.EndpointRule().Update(txctx, epr)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleUpdateRes{Doc: epr.(*entities.EndpointRule)}
	return res, nil
}
