package repositories

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	Create(ctx context.Context, ws *entities.Application) (*entities.Application, error)
	Get(ctx context.Context, id string) (*entities.Application, error)
	List(ctx context.Context, wsId, name string) ([]entities.Application, error)
	Update(ctx context.Context, ws *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, id string) (*entities.Application, error)

	GetWithWorkspace(ctx context.Context, id string) (*ApplicationWithWorkspace, error)
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
