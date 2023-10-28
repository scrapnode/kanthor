package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointListReq struct {
	WsId  string
	AppId string
	*structure.ListReq
}

func (req *EndpointListReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.PointerNotNil("list", req.ListReq),
	)
}

type EndpointListRes struct {
	*structure.ListRes[entities.Endpoint]
}

func (uc *endpoint) List(ctx context.Context, req *EndpointListReq) (*EndpointListRes, error) {
	app, err := uc.repos.Application().Get(ctx, req.WsId, req.AppId)
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
