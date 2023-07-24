package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Endpoint interface {
	Create(ctx context.Context, entity *entities.Endpoint) (*entities.Endpoint, error)
	BulkCreate(ctx context.Context, entities []entities.Endpoint) ([]string, error)

	List(ctx context.Context, wsId, appId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error)
	Get(ctx context.Context, wsId, appId, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, entity *entities.EndpointRule) (*entities.EndpointRule, error)
	BulkCreate(ctx context.Context, entities []entities.EndpointRule) ([]string, error)

	List(ctx context.Context, wsId, appId, epId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error)
}
