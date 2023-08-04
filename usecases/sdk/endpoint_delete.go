package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *endpoint) Delete(ctx context.Context, req *EndpointDeleteReq) (*EndpointDeleteRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Endpoint().Get(txctx, ws.Id, req.AppId, req.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repos.Endpoint().Delete(txctx, app); err != nil {
			return nil, err
		}
		return app, nil
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointDeleteRes{Doc: app.(*entities.Endpoint)}
	return res, nil
}
