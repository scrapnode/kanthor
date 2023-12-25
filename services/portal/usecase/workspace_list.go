package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
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
	Data  []entities.Workspace
	Count int64
}

func (uc *workspace) List(ctx context.Context, in *WorkspaceListIn) (*WorkspaceListOut, error) {
	out := &WorkspaceListOut{}

	// owner
	own, err := uc.repositories.Workspace().ListOwned(ctx, in.AccId)
	if err != nil {
		return nil, err
	}
	out.Data = append(out.Data, own...)

	// collaborator
	tenants, err := uc.infra.Authorizator.Tenants(in.AccId)
	if err != nil {
		return nil, err
	}
	if len(tenants) > 0 {
		workspaces, err := uc.repositories.Workspace().ListByIds(ctx, tenants)
		if err != nil {
			return nil, err
		}
		out.Data = append(out.Data, workspaces...)
	}

	return out, nil
}
