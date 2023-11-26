package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceUpdateIn struct {
	Id   string
	Name string
}

func (in *WorkspaceUpdateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
		validator.StringRequired("name", in.Name),
	)
}

type WorkspaceUpdateOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error) {
	ws, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repositories.Workspace().Get(txctx, in.Id)
		if err != nil {
			return nil, err
		}

		ws.Name = in.Name
		ws.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Workspace().Update(txctx, ws)
	})
	if err != nil {
		return nil, err
	}

	return &WorkspaceUpdateOut{Doc: ws.(*entities.Workspace)}, nil
}
