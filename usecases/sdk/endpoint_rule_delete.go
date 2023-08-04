package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *endpointRule) Delete(ctx context.Context, req *EndpointRuleDeleteReq) (*EndpointRuleDeleteRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.EndpointRule().Get(txctx, ws.Id, req.AppId, req.EpId, req.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repos.EndpointRule().Delete(txctx, app); err != nil {
			return nil, err
		}
		return app, nil
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleDeleteRes{Doc: app.(*entities.EndpointRule)}
	return res, nil
}
