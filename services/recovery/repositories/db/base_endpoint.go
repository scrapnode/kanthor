package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Endpoint interface {
	List(ctx context.Context, appId string) ([]entities.Endpoint, error)
	Rules(ctx context.Context, appId string) ([]entities.EndpointRule, error)
}
