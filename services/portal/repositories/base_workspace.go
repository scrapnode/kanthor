package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/structure"
)

type Workspace interface {
	Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	GetOwned(ctx context.Context, owner string) (*entities.Workspace, error)
	ListOwned(ctx context.Context, owner string) ([]entities.Workspace, error)
}

type WorkspaceCredentials interface {
	Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
	Update(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
	Get(ctx context.Context, wsId, id string) (*entities.WorkspaceCredentials, error)
	List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.WorkspaceCredentials], error)
}
