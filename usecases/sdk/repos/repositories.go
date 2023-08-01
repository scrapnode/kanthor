package repos

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *database.Config, logger logging.Logger) Repositories {
	return NewSql(conf, logger)
}

type Repositories interface {
	patterns.Connectable
	Transaction(ctx context.Context, handler func(txctx context.Context) (interface{}, error)) (res interface{}, err error)
	Workspace() Workspace
	WorkspaceTier() WorkspaceTier
	WorkspaceCredentials() WorkspaceCredentials
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
}
