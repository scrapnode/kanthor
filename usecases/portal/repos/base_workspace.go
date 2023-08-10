package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Workspace interface {
	BulkCreate(ctx context.Context, docs []entities.Workspace) ([]string, error)

	Get(ctx context.Context, id string) (*entities.Workspace, error)
	GetOwned(ctx context.Context, owner string) (*entities.Workspace, error)
}

type WorkspaceTier interface {
	BulkCreate(ctx context.Context, docs []entities.WorkspaceTier) ([]string, error)

	Create(ctx context.Context, doc *entities.WorkspaceTier) (*entities.WorkspaceTier, error)
	Get(ctx context.Context, wsId string) (*entities.WorkspaceTier, error)
}

type WorkspaceCredentials interface {
	BulkCreate(ctx context.Context, docs []entities.WorkspaceCredentials) ([]string, error)

	Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
}
