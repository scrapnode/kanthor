package repos

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Application interface {
	Create(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Update(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, doc *entities.Application) error

	List(ctx context.Context, ws *entities.Workspace, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error)
	Get(ctx context.Context, ws *entities.Workspace, id string) (*entities.Application, error)
}
