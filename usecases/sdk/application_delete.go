package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *application) Delete(ctx context.Context, req *ApplicationDeleteReq) (*ApplicationDeleteRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Application().Get(txctx, ws.Id, req.Id)
		if err != nil {
			return nil, err
		}

		if err := uc.repos.Application().Delete(txctx, app); err != nil {
			return nil, err
		}
		return app, nil
	})
	if err != nil {
		return nil, err
	}

	res := &ApplicationDeleteRes{Doc: app.(*entities.Application)}
	return res, nil
}
