package portal

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
)

func (uc *workspaceCredentials) Get(ctx context.Context, req *WorkspaceCredentialsGetReq) (*WorkspaceCredentialsGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	// we don't need to use cache here because the usage is too low
	wsc, err := uc.repos.WorkspaceCredentials().Get(ctx, ws.Id, req.Id)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: don't return hash value
	wsc.Hash = ""

	res := &WorkspaceCredentialsGetRes{Doc: wsc}
	return res, nil
}
