package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/project"
)

type WorkspaceCreateIn struct {
	AccId string
	Name  string
}

func (in *WorkspaceCreateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("acc_id", in.AccId),
		validator.StringRequired("name", in.Name),
	)
}

type WorkspaceCreateOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Create(ctx context.Context, in *WorkspaceCreateIn) (*WorkspaceCreateOut, error) {
	res, err := uc.repositories.Database().Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		doc := &entities.Workspace{
			OwnerId: in.AccId,
			Name:    in.Name,
			Tier:    project.Tier(),
		}
		doc.Id = identifier.New(entities.IdNsWs)
		doc.SetAT(uc.infra.Timer.Now())

		ws, err := uc.repositories.Database().Workspace().Create(ctx, doc)
		if err != nil {
			return nil, err
		}

		return &WorkspaceCreateOut{Doc: ws}, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*WorkspaceCreateOut), nil
}
