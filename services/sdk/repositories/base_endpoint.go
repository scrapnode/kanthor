package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Endpoint interface {
	Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, doc *entities.Endpoint) error

	List(ctx context.Context, wsId, appId string, q string, limit, page int) ([]entities.Endpoint, error)
	Count(ctx context.Context, wsId, appId string, q string) (int64, error)
	Get(ctx context.Context, wsId string, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Update(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, doc *entities.EndpointRule) error

	List(ctx context.Context, wsId, epId string, q string, limit, page int) ([]entities.EndpointRule, error)
	Count(ctx context.Context, wsId, epId string, q string) (int64, error)
	Get(ctx context.Context, wsId string, id string) (*entities.EndpointRule, error)
}
