package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleDeleteReq struct {
	EpId string
	Id   string
}

func (req *EndpointRuleDeleteReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", req.EpId, entities.IdNsEp),
		validator.StringStartsWith("id", req.EpId, entities.IdNsEpr),
	)
}

type EndpointRuleDeleteRes struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Delete(ctx context.Context, req *EndpointRuleDeleteReq) (*EndpointRuleDeleteRes, error) {
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

		if err := uc.repos.EndpointRule().Delete(txctx, epr); err != nil {
			return nil, err
		}
		return epr, nil
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleDeleteRes{Doc: epr.(*entities.EndpointRule)}
	return res, nil
}
