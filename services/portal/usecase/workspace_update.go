package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceUpdateIn struct {
	AccId string
	Id    string
	Name  string
}

func (in *WorkspaceUpdateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("acc_id", in.AccId),
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
		validator.StringRequired("name", in.Name),
	)
}

type WorkspaceUpdateOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error) {
	getout, err := uc.Get(ctx, &WorkspaceGetIn{AccId: in.AccId, Id: in.Id})
	if err != nil {
		return nil, err
	}

	ws, err := uc.repositories.Database().Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repositories.Database().Workspace().Get(txctx, getout.Doc.Id)
		if err != nil {
			return nil, err
		}

		ws.Name = in.Name
		ws.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Database().Workspace().Update(txctx, ws)
	})
	if err != nil {
		return nil, err
	}

	return &WorkspaceUpdateOut{Doc: ws.(*entities.Workspace)}, nil
}
