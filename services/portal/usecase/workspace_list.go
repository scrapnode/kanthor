package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceListIn struct {
	AccId string
}

func (in *WorkspaceListIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("acc_id", in.AccId),
	)
}

type WorkspaceListOut struct {
	Workspaces map[string]*entities.Workspace
}

func (uc *workspace) List(ctx context.Context, in *WorkspaceListIn) (*WorkspaceListOut, error) {
	out := &WorkspaceListOut{Workspaces: map[string]*entities.Workspace{}}

	// owner
	own, err := uc.repositories.Workspace().ListOwned(ctx, in.AccId)
	if err != nil {
		return nil, err
	}
	for _, ws := range own {
		out.Workspaces[ws.Id] = &ws
	}

	// collaborator
	tenants, err := uc.infra.Authorizator.Tenants(in.AccId)
	if err != nil {
		return nil, err
	}
	if len(tenants) > 0 {
		cooperates, err := uc.repositories.Workspace().List(ctx, structure.WithListIds(tenants))
		if err != nil {
			return nil, err
		}
		for _, ws := range cooperates.Data {
			out.Workspaces[ws.Id] = &ws
		}

	}

	return out, nil
}
