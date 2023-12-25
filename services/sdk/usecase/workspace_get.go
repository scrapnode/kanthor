package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceGetIn struct {
	Id string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	ws, err := uc.repositories.Workspace().Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &WorkspaceGetOut{Doc: ws}, nil
}
