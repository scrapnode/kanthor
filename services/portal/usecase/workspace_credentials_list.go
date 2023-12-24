package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsListIn struct {
	*entities.Query
	WsId string
}

func (in *WorkspaceCredentialsListIn) Validate() error {
	if err := in.Query.Validate(); err != nil {
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
	data, err := uc.repositories.WorkspaceCredentials().List(ctx, in.WsId, in.Query)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.WorkspaceCredentials().Count(ctx, in.WsId, in.Query)
	if err != nil {
		return nil, err
	}

	out := &WorkspaceCredentialsListOut{Data: data, Count: count}
	return out, nil
}
