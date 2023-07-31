package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *application) Create(ctx context.Context, req *ApplicationCreateReq) (*ApplicationCreateRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	doc := &entities.Application{WorkspaceId: ws.Id, Name: req.Name}
	doc.GenId()
	doc.SetAT(uc.timer.Now())

	app, err := uc.repos.Application().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &ApplicationCreateRes{Doc: app}
	return res, nil
}
