package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Endpoint interface {
	Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, doc *entities.Endpoint) error

	List(ctx context.Context, wsId, appId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error)
	Get(ctx context.Context, wsId, appId, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Update(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, doc *entities.EndpointRule) error

	List(ctx context.Context, wsId, appId, epId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error)
	Get(ctx context.Context, wsId, appId, epId, id string) (*entities.EndpointRule, error)
}
