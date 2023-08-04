package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *endpoint) Update(ctx context.Context, req *EndpointUpdateReq) (*EndpointUpdateRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Endpoint().Get(txctx, ws.Id, req.AppId, req.Id)
		if err != nil {
			return nil, err
		}

		app.Name = req.Name
		app.SetAT(uc.timer.Now())
		return uc.repos.Endpoint().Update(txctx, app)
	})
	if err != nil {
		return nil, err
	}

	res := &EndpointUpdateRes{Doc: app.(*entities.Endpoint)}
	return res, nil
}
