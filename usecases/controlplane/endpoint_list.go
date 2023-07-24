package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/structure"
)

func (uc *endpoint) List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error) {
	list, err := uc.repos.Endpoint().List(
		ctx,
		req.Workspace.Id,
		req.AppId,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &EndpointListRes{Cursor: list.Cursor, Data: list.Data}
	return res, nil
}
