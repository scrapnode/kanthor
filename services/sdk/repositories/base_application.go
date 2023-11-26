package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/structure"
)

type Application interface {
	Create(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Update(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, doc *entities.Application) error

	List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
