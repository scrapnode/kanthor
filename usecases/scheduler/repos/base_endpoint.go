package repos

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Endpoint interface {
	List(ctx context.Context, appId string) ([]entities.Endpoint, error)
	Rules(ctx context.Context, appId string) ([]entities.EndpointRule, error)
}
