package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsListIn struct {
	WsId string
	*structure.ListReq
}

func (in *WorkspaceCredentialsListIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.PointerNotNil("list", in.ListReq),
	)
}

type WorkspaceCredentialsListOut struct {
	*structure.ListRes[entities.WorkspaceCredentials]
}

func (uc *workspaceCredentials) List(ctx context.Context, in *WorkspaceCredentialsListIn) (*WorkspaceCredentialsListOut, error) {
	listing, err := uc.repositories.WorkspaceCredentials().List(
		ctx, in.WsId,
		structure.WithListCursor(in.Cursor),
		structure.WithListSearch(in.Search),
		structure.WithListLimit(in.Limit),
		structure.WithListIds(in.Ids),
	)
	if err != nil {
		return nil, err
	}

	for i, wsc := range listing.Data {
		// IMPORTANT: don't return hash value
		wsc.Hash = ""
		listing.Data[i] = wsc
	}
	return &WorkspaceCredentialsListOut{ListRes: listing}, nil
}
