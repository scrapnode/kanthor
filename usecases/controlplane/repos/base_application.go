package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Application interface {
	Create(ctx context.Context, entity *entities.Application) (*entities.Application, error)
	BulkCreate(ctx context.Context, entities []entities.Application) ([]string, error)

	List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
