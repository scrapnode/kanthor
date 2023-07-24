package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/structure"
)

func (uc *endpointRule) List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error) {
	list, err := uc.repos.EndpointRule().List(
		ctx,
		req.Workspace.Id,
		req.AppId,
		req.EpId,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleListRes{Cursor: list.Cursor, Data: list.Data}
	return res, nil
}
