package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointListIn struct {
	WsId  string
	AppId string
	*structure.ListReq
}

func (in *EndpointListIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.PointerNotNil("list", in.ListReq),
	)
}

type EndpointListOut struct {
	*structure.ListRes[entities.Endpoint]
}

func (uc *endpoint) List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error) {
	app, err := uc.repositories.Application().Get(ctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	listing, err := uc.repositories.Endpoint().List(
		ctx,
		app,
		structure.WithListCursor(in.Cursor),
		structure.WithListSearch(in.Search),
		structure.WithListLimit(in.Limit),
		structure.WithListIds(in.Ids),
	)
	if err != nil {
		return nil, err
	}

	return &EndpointListOut{ListRes: listing}, nil
}
