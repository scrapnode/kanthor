package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Workspace interface {
	Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	ListByIds(ctx context.Context, ids []string) ([]entities.Workspace, error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	GetOwned(ctx context.Context, owner, id string) (*entities.Workspace, error)
	ListOwned(ctx context.Context, owner string) ([]entities.Workspace, error)
	GetSnapshotRows(ct context.Context, id string) ([]SnaptshotRow, error)
}

type WorkspaceCredentials interface {
	Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
	Update(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error)
	List(ctx context.Context, wsId string, query *entities.PagingQuery) ([]entities.WorkspaceCredentials, error)
	Count(ctx context.Context, wsId string, query *entities.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.WorkspaceCredentials, error)
}
