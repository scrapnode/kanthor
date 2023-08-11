package portal

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
)

func (uc *workspaceCredentials) Update(ctx context.Context, req *WorkspaceCredentialsUpdateReq) (*WorkspaceCredentialsUpdateRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	doc, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsc, err := uc.repos.WorkspaceCredentials().Get(txctx, ws.Id, req.Id)
		if err != nil {
			return nil, err
		}

		wsc.Name = req.Name
		wsc.SetAT(uc.timer.Now())
		return uc.repos.WorkspaceCredentials().Update(txctx, wsc)
	})
	if err != nil {
		return nil, err
	}

	wsc := doc.(*entities.WorkspaceCredentials)
	// IMPORTANT: don't return hash value
	wsc.Hash = ""

	res := &WorkspaceCredentialsUpdateRes{Doc: wsc}
	return res, nil
}
