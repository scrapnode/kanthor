package repositories

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Workspace interface {
	Create(ctx context.Context, ws *entities.Workspace) (*entities.Workspace, error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	List(ctx context.Context, opts ...ListOps) (*ListRes[entities.Workspace], error)
	Update(ctx context.Context, ws *entities.Workspace) (*entities.Workspace, error)
	Delete(ctx context.Context, id string) (*entities.Workspace, error)
}
