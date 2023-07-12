package repositories

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Endpoint interface {
	Create(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error)
	Get(ctx context.Context, id string) (*entities.Endpoint, error)
	List(ctx context.Context, appId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error)
	Update(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, id string) (*entities.Endpoint, error)
}
