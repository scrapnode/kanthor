package repositories

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

func New(conf *database.Config, logger logging.Logger, timer timer.Timer) Repositories {
	return NewSql(conf, logger, timer)
}

type Repositories interface {
	patterns.Connectable
	Workspace() Workspace
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
}

type Workspace interface {
	Create(ctx context.Context, ws *entities.Workspace) (*entities.Workspace, error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	List(ctx context.Context, name string) ([]entities.Workspace, error)
	Update(ctx context.Context, ws *entities.Workspace) (*entities.Workspace, error)
	Delete(ctx context.Context, id string) (*entities.Workspace, error)
}

type Application interface {
	Create(ctx context.Context, ws *entities.Application) (*entities.Application, error)
	Get(ctx context.Context, id string) (*entities.Application, error)
	List(ctx context.Context, wsId, name string) ([]entities.Application, error)
	Update(ctx context.Context, ws *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, id string) (*entities.Application, error)
}

type Endpoint interface {
	Create(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error)
	Get(ctx context.Context, id string) (*entities.Endpoint, error)
	List(ctx context.Context, appId, name string) ([]entities.Endpoint, error)
	Update(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error)
	Get(ctx context.Context, id string) (*entities.EndpointRule, error)
	List(ctx context.Context, epId string) ([]entities.EndpointRule, error)
	Update(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, id string) (*entities.EndpointRule, error)
}
