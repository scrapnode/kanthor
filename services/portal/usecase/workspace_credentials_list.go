package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsListIn struct {
	*entities.PagingQuery
	WsId string
}

func (in *WorkspaceCredentialsListIn) Validate() error {
	if err := in.PagingQuery.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
	)
}

type WorkspaceCredentialsListOut struct {
	Data  []entities.WorkspaceCredentials
	Count int64
}

func (uc *workspaceCredentials) List(ctx context.Context, in *WorkspaceCredentialsListIn) (*WorkspaceCredentialsListOut, error) {
	data, err := uc.repositories.Database().WorkspaceCredentials().List(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Database().WorkspaceCredentials().Count(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	out := &WorkspaceCredentialsListOut{Data: data, Count: count}
	return out, nil
}
