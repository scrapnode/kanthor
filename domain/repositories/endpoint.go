package repositories

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Endpoint interface {
	Create(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error)
	Get(ctx context.Context, id string) (*entities.Endpoint, error)
	List(ctx context.Context, appId, name string) ([]entities.Endpoint, error)
	Update(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error)
	Get(ctx context.Context, id string) (*entities.EndpointRule, error)
	List(ctx context.Context, epId string) ([]entities.EndpointRule, error)
	Update(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, id string) (*entities.EndpointRule, error)
}
