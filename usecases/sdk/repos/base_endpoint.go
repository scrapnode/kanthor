package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Endpoint interface {
	Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Update(ctx context.Context, wsId string, doc *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, wsId, id string) error

	List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error)
	Get(ctx context.Context, wsId, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Update(ctx context.Context, wsId string, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, wsId, id string) error

	List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error)
	Get(ctx context.Context, wsId, id string) (*entities.EndpointRule, error)
}
