package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Endpoint interface {
	BulkCreate(ctx context.Context, docs []entities.Endpoint) ([]string, error)
}

type EndpointRule interface {
	BulkCreate(ctx context.Context, docs []entities.EndpointRule) ([]string, error)
}
