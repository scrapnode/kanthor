package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationListReq struct {
	WorkspaceId string
	*structure.ListReq
}

func (req *ApplicationListReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WorkspaceId, entities.IdNsWs),
		validator.PointerNotNil("list", req.ListReq),
	)
}

type ApplicationListRes struct {
	*structure.ListRes[entities.Application]
}

func (uc *application) List(ctx context.Context, req *ApplicationListReq) (*ApplicationListRes, error) {
	listing, err := uc.repos.Application().List(
		ctx,
		req.WorkspaceId,
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
