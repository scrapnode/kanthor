package rest

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

func RegisterWorkspaceResolver(uc usecase.Portal) func(ctx context.Context, id string) (*entities.Workspace, error) {
	return func(ctx context.Context, id string) (*entities.Workspace, error) {
		in := &usecase.WorkspaceGetIn{Id: id}
		if err := in.Validate(); err != nil {
			return nil, err
		}

		out, err := uc.Workspace().Get(ctx, in)
		if err != nil {
			return nil, err
		}
		return out.Doc, nil
	}
}
