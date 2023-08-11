package repos

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(conf *database.Config, logger logging.Logger, timer timer.Timer) Repositories {
	return NewSql(conf, logger, timer)
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
