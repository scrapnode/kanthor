package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceGetIn struct {
	AccId string
	Id    string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("acc_id", in.AccId),
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	ws, err := uc.repositories.Workspace().GetOwned(ctx, in.AccId, in.Id)
	if err == nil {
		return &WorkspaceGetOut{Doc: ws}, nil
	}

	tenants, err := uc.infra.Authorizator.Tenants(in.AccId)
	if err == nil && slices.Contains(tenants, in.Id) {
		ws, err := uc.repositories.Workspace().Get(ctx, in.Id)
		if err == nil {
			return &WorkspaceGetOut{Doc: ws}, nil
		}
	}

	return nil, errors.New("you don't own the workspace or have no permission to access it")
}
