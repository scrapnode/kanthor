package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Application interface {
	Create(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Update(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, doc *entities.Application) error

	List(ctx context.Context, wsId string, q string, limit, page int) ([]entities.Application, error)
	Count(ctx context.Context, wsId string, q string) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
