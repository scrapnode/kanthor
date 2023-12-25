package rest

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

func RegisterWorkspaceResolver(uc usecase.Portal) func(ctx context.Context, acc *authenticator.Account, id string) (*entities.Workspace, error) {
	return func(ctx context.Context, acc *authenticator.Account, id string) (*entities.Workspace, error) {
		in := &usecase.WorkspaceGetIn{
			AccId: acc.Sub,
			Id:    id,
		}
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
