package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointListReq struct {
	AppId string
	*structure.ListReq
}

func (req *EndpointListReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, "app_"),
		validator.PointerNotNil("list", req.ListReq),
	)
}

type EndpointListRes struct {
	*structure.ListRes[entities.Endpoint]
}

func (uc *endpoint) List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	app, err := uc.repos.Application().Get(ctx, ws, req.AppId)
	if err != nil {
		return nil, err
	}

	listing, err := uc.repos.Endpoint().List(
		ctx,
		app,
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
