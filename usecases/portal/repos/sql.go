package repos

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

func NewSql(conf *database.Config, logger logging.Logger, timer timer.Timer) Repositories {
	db := database.NewSQL(conf, logger)

	logger = logger.With("component", "repositories.sql")
	return &sql{logger: logger, timer: timer, db: db}
}

type sql struct {
	logger logging.Logger
	timer  timer.Timer
	db     database.Database

	client               *gorm.DB
	workspace            *SqlWorkspace
	workspaceTier        *SqlWorkspaceTier
	workspaceCredentials *SqlWorkspaceCredentials
	application          *SqlApplication
	endpoint             *SqlEndpoint
	endpointRule         *SqlEndpointRule
}

func (repos *sql) Connect(ctx context.Context) error {
	if err := repos.db.Connect(ctx); err != nil {
		return err
	}

	repos.client = repos.db.Client().(*gorm.DB)
	repos.logger.Info("connected")
	return nil
}

func (repos *sql) Disconnect(ctx context.Context) error {
	repos.logger.Info("disconnected")

	if err := repos.db.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (repos *sql) Transaction(ctx context.Context, handler func(txctx context.Context) (interface{}, error)) (res interface{}, err error) {
	err = repos.client.Transaction(func(tx *gorm.DB) error {
		res, err = handler(context.WithValue(ctx, database.CtxTransaction, tx))
		return err
	})
	return
}

func (repos *sql) Workspace() Workspace {
	if repos.workspace == nil {
		repos.workspace = &SqlWorkspace{client: repos.client, timer: repos.timer}
	}

	return repos.workspace
}

func (repos *sql) WorkspaceTier() WorkspaceTier {
	if repos.workspaceTier == nil {
		repos.workspaceTier = &SqlWorkspaceTier{client: repos.client, timer: repos.timer}
	}

	return repos.workspaceTier
}

func (repos *sql) WorkspaceCredentials() WorkspaceCredentials {
	if repos.workspaceCredentials == nil {
		repos.workspaceCredentials = &SqlWorkspaceCredentials{client: repos.client, timer: repos.timer}
	}

	return repos.workspaceCredentials
}

func (repos *sql) Application() Application {
	if repos.application == nil {
		repos.application = &SqlApplication{client: repos.client, timer: repos.timer}
	}

	return repos.application
}

func (repos *sql) Endpoint() Endpoint {
	if repos.endpoint == nil {
		repos.endpoint = &SqlEndpoint{client: repos.client, timer: repos.timer}
	}

	return repos.endpoint
}

func (repos *sql) EndpointRule() EndpointRule {
	if repos.endpointRule == nil {
		repos.endpointRule = &SqlEndpointRule{client: repos.client, timer: repos.timer}
	}

	return repos.endpointRule
}
