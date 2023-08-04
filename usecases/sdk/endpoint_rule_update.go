package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *endpointRule) Update(ctx context.Context, req *EndpointRuleUpdateReq) (*EndpointRuleUpdateRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.EndpointRule().Get(txctx, ws.Id, req.AppId, req.EpId, req.Id)
		if err != nil {
			return nil, err
		}

		app.Name = req.Name
		app.SetAT(uc.timer.Now())
		return uc.repos.EndpointRule().Update(txctx, app)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleUpdateRes{Doc: app.(*entities.EndpointRule)}
	return res, nil
}
