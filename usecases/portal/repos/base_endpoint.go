package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Endpoint interface {
	BulkCreate(ctx context.Context, docs []entities.Endpoint) ([]string, error)
}

type EndpointRule interface {
	BulkCreate(ctx context.Context, docs []entities.EndpointRule) ([]string, error)
}
