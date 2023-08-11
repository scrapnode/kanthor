package portal

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
)

func (uc *workspaceCredentials) List(ctx context.Context, req *WorkspaceCredentialsListReq) (*WorkspaceCredentialsListRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	listing, err := uc.repos.WorkspaceCredentials().List(
		ctx, ws.Id,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	for i, wsc := range listing.Data {
		// IMPORTANT: don't return hash value
		wsc.Hash = ""
		listing.Data[i] = wsc
	}
	res := &WorkspaceCredentialsListRes{ListRes: listing}
	return res, nil
}
