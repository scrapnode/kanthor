package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
)

func (uc *application) List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	listing, err := uc.repos.Application().List(
		ctx, ws.Id,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &ApplicationListRes{ListRes: listing}
	return res, nil
}
