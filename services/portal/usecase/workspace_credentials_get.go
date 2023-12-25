package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsGetIn struct {
	WsId string
	Id   string
}

func (in *WorkspaceCredentialsGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsWsc),
	)
}

type WorkspaceCredentialsGetOut struct {
	Doc *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Get(ctx context.Context, in *WorkspaceCredentialsGetIn) (*WorkspaceCredentialsGetOut, error) {
	// we don't need to use cache here because the usage is too low
	wsc, err := uc.repositories.Database().WorkspaceCredentials().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: don't return hash value
	wsc.Hash = ""

	return &WorkspaceCredentialsGetOut{Doc: wsc}, nil
}
