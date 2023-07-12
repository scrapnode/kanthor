package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type EndpointRule interface {
	Create(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error)
	Get(ctx context.Context, id string) (*entities.EndpointRule, error)
	List(ctx context.Context, epId string) (*structure.ListRes[entities.EndpointRule], error)
	Update(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, id string) (*entities.EndpointRule, error)
}
