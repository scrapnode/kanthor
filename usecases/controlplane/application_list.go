package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/structure"
)

func (uc *application) List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error) {
	list, err := uc.repos.Application().List(
		ctx,
		req.Workspace.Id,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &ApplicationListRes{Cursor: list.Cursor, Data: list.Data}
	return res, nil
}
