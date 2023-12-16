package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationListIn struct {
	WsId string
	*structure.ListReq
}

func (in *ApplicationListIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.PointerNotNil("list", in.ListReq),
	)
}

type ApplicationListOut struct {
	*structure.ListRes[entities.Application]
}

func (uc *application) List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error) {
	listing, err := uc.repositories.Application().List(
		ctx,
		in.WsId,
		structure.WithListCursor(in.Cursor),
		structure.WithListSearch(in.Search),
		structure.WithListLimit(in.Limit),
		structure.WithListIds(in.Ids),
	)
	if err != nil {
		return nil, err
	}

	return &ApplicationListOut{ListRes: listing}, nil
}
