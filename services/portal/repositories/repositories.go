package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
)

func New(logger logging.Logger, db database.Database) Repositories {
	return NewSql(logger, db)
}

type Repositories interface {
	Transaction(ctx context.Context, handler func(txctx context.Context) (interface{}, error)) (res interface{}, err error)
	Workspace() Workspace
	WorkspaceCredentials() WorkspaceCredentials
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
}
