package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	Get(ctx context.Context, id string) (*entities.Application, error)
	// ListEndpointsWithRules return list of endpoints with their active rules. They are well sorted list with logic:
	//
	// rule.priority - rule.exclusionary
	//			  15 - TRUE
	//			  15 - FALSE
	//			  9  - FALSE
	//			  70 - TRUE
	//			  70 - FALSE
	//			  0  - FALSE
	//
	// IMPORTANT: the order of the list above is important
	ListEndpointsWithRules(ctx context.Context, id string) (*ApplicationWithEndpointsAndRules, error)
}

type ApplicationWithEndpointsAndRules struct {
	entities.Application
	Endpoints []EndpointWithRules
}

type EndpointWithRules struct {
	entities.Endpoint
	Rules []entities.EndpointRule
}
