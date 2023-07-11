package repositories

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	Create(ctx context.Context, ws *entities.Application) (*entities.Application, error)
	Get(ctx context.Context, id string) (*entities.Application, error)
	List(ctx context.Context, wsId string, opts ...ListOps) (*ListRes[entities.Application], error)
	Update(ctx context.Context, ws *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, id string) (*entities.Application, error)

	GetWithWorkspace(ctx context.Context, id string) (*ApplicationWithWorkspace, error)

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

type ApplicationWithWorkspace struct {
	entities.Application
	Workspace entities.Workspace
}

type ApplicationWithEndpointsAndRules struct {
	entities.Application
	Endpoints []EndpointWithRules
}

type EndpointWithRules struct {
	entities.Endpoint
	Rules []entities.EndpointRule
}
