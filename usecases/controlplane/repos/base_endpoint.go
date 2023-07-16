package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Endpoint interface {
	Create(ctx context.Context, entity *entities.Endpoint) (*entities.Endpoint, error)
	BulkCreate(ctx context.Context, entities []entities.Endpoint) ([]string, error)
}

type EndpointRule interface {
	Create(ctx context.Context, entity *entities.EndpointRule) (*entities.EndpointRule, error)
	BulkCreate(ctx context.Context, entities []entities.EndpointRule) ([]string, error)
}
