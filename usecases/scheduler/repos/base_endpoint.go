package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Endpoint interface {
	// List function of Endpoint is not verified by workspace id
	// because it will be used in a safe internal place
	// where we all verified request before execute the query
	List(ctx context.Context, appId string) ([]entities.Endpoint, error)
}

type EndpointRule interface {
	// List function of EndpointRule is not verified by workspace id
	// because it will be used in a safe internal place
	// where we all verified request before execute the query
	List(ctx context.Context, epIds []string) ([]entities.EndpointRule, error)
}
