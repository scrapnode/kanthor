package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Application interface {
	BulkCreate(ctx context.Context, docs []entities.Application) ([]string, error)

	Create(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
