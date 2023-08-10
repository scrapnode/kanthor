package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
)

func (uc *endpointRule) List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	listing, err := uc.repos.EndpointRule().List(
		ctx, ws.Id, req.AppId, req.EpId,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleListRes{ListRes: listing}
	return res, nil
}
