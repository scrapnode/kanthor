package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Endpoint interface {
	BulkCreate(ctx context.Context, docs []entities.Endpoint) ([]string, error)
	Count(ctx context.Context, wsId string, query *entities.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	BulkCreate(ctx context.Context, docs []entities.EndpointRule) ([]string, error)
}
