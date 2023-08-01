package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *application) Update(ctx context.Context, req *ApplicationUpdateReq) (*ApplicationUpdateRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Application().Get(txctx, ws.Id, req.Id)
		if err != nil {
			return nil, err
		}

		app.Name = req.Name
		app.SetAT(uc.timer.Now())
		return uc.repos.Application().Update(txctx, app)
	})
	if err != nil {
		return nil, err
	}

	res := &ApplicationUpdateRes{Doc: app.(*entities.Application)}
	return res, nil
}
