package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointUpdateReq struct {
	AppId string
	Id    string
	Name  string
}

func (req *EndpointUpdateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", req.Id, entities.IdNsEp),
		validator.StringRequired("name", req.Name),
	)
}

type EndpointUpdateRes struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Update(ctx context.Context, req *EndpointUpdateReq) (*EndpointUpdateRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	ep, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Application().Get(ctx, ws, req.AppId)
		if err != nil {
			return nil, err
		}

		ep, err := uc.repos.Endpoint().Get(txctx, app, req.Id)
		if err != nil {
			return nil, err
		}

		ep.Name = req.Name
		ep.SetAT(uc.infra.Timer.Now())
		return uc.repos.Endpoint().Update(txctx, ep)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointUpdateRes{Doc: ep.(*entities.Endpoint)}
	return res, nil
}
