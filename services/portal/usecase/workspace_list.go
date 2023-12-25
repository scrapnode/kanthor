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
	seen := map[string]bool{}

	// owner
	own, err := uc.repositories.Workspace().ListOwned(ctx, in.AccId)
	if err != nil {
		return nil, err
	}
	for _, ws := range own {
		if _, found := seen[ws.Id]; found {
			continue
		}

		seen[ws.Id] = true
		out.Data = append(out.Data, ws)
	}

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
		for _, ws := range workspaces {
			if _, found := seen[ws.Id]; found {
				continue
			}

			seen[ws.Id] = true
			out.Data = append(out.Data, ws)
		}
	}

	return out, nil
}
