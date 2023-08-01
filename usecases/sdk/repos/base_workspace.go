package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
}

type WorkspaceTier interface {
	Get(ctx context.Context, wsId string) (*entities.WorkspaceTier, error)
}

type WorkspaceCredentials interface {
	Get(ctx context.Context, id string) (*entities.WorkspaceCredentials, error)
}
