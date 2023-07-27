package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Workspace interface {
	BulkCreate(ctx context.Context, docs []entities.Workspace) ([]string, error)

	Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	GetOwned(ctx context.Context, owner string) (*entities.Workspace, error)
}

type WorkspaceTier interface {
	BulkCreate(ctx context.Context, docs []entities.WorkspaceTier) ([]string, error)

	Create(ctx context.Context, doc *entities.WorkspaceTier) (*entities.WorkspaceTier, error)
}

type WorkspaceCredentials interface {
	BulkCreate(ctx context.Context, docs []entities.WorkspaceCredentials) ([]string, error)

	Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
}