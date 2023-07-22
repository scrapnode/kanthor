package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	Get(ctx context.Context, id string) (*entities.Application, error)
	Create(ctx context.Context, entity *entities.Application) (*entities.Application, error)
	BulkCreate(ctx context.Context, entities []entities.Application) ([]string, error)
}
