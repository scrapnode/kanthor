package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Workspace interface {
	Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	ListByIds(ctx context.Context, ids []string) (*[]entities.Workspace, error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	GetOwned(ctx context.Context, owner string) (*entities.Workspace, error)
	ListOwned(ctx context.Context, owner string) ([]entities.Workspace, error)
}

type WorkspaceCredentials interface {
	Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
	Update(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
	List(ctx context.Context, wsId string, limit, page int, q string) ([]entities.WorkspaceCredentials, error)
	Count(ctx context.Context, wsId string, q string) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.WorkspaceCredentials, error)
}
