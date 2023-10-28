package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsListReq struct {
	WsId string
	*structure.ListReq
}

func (req *WorkspaceCredentialsListReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.PointerNotNil("list", req.ListReq),
	)
}

type WorkspaceCredentialsListRes struct {
	*structure.ListRes[entities.WorkspaceCredentials]
}

func (uc *workspaceCredentials) List(ctx context.Context, req *WorkspaceCredentialsListReq) (*WorkspaceCredentialsListRes, error) {
	listing, err := uc.repositories.WorkspaceCredentials().List(
		ctx, req.WsId,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	for i, wsc := range listing.Data {
		// IMPORTANT: don't return hash value
		wsc.Hash = ""
		listing.Data[i] = wsc
	}
	res := &WorkspaceCredentialsListRes{ListRes: listing}
	return res, nil
}
