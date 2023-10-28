package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
}

type WorkspaceCredentials interface {
	Get(ctx context.Context, id string) (*entities.WorkspaceCredentials, error)
}
