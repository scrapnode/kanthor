package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
)

func (uc *endpoint) List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	listing, err := uc.repos.Endpoint().List(
		ctx, ws.Id, req.AppId,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &EndpointListRes{ListRes: listing}
	return res, nil
}
