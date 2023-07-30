package repos

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
	"sync"
)

func NewSql(conf *database.Config, logger logging.Logger, timer timer.Timer) Repositories {
	db := database.NewSQL(conf, logger)

	logger = logger.With("component", "repoitories.sql")
	return &sql{logger: logger, timer: timer, db: db}
}

type sql struct {
	logger logging.Logger
	timer  timer.Timer
	db     database.Database

	mu                   sync.RWMutex
	client               *gorm.DB
	workspace            *SqlWorkspace
	workspaceTier        *SqlWorkspaceTier
	workspaceCredentials *SqlWorkspaceCredentials
	application          *SqlApplication
	endpoint             *SqlEndpoint
	endpointRule         *SqlEndpointRule
}

func (repo *sql) Connect(ctx context.Context) error {
	if err := repo.db.Connect(ctx); err != nil {
		return err
	}

	repo.client = repo.db.Client().(*gorm.DB)
	repo.logger.Info("connected")
	return nil
}

func (repo *sql) Disconnect(ctx context.Context) error {
	repo.logger.Info("disconnected")

	if err := repo.db.Disconnect(ctx); err != nil {
		return err
	}

	repo.client = nil

	return nil
}

func (repo *sql) Transaction(ctx context.Context, handler func(txctx context.Context) (interface{}, error)) (res interface{}, err error) {
	err = repo.client.Transaction(func(tx *gorm.DB) error {
		res, err = handler(context.WithValue(ctx, database.CtxTransaction, tx))
		return err
	})
	return
}

func (repo *sql) Workspace() Workspace {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspace == nil {
		repo.workspace = &SqlWorkspace{client: repo.client, timer: repo.timer}
	}

	return repo.workspace
}

func (repo *sql) WorkspaceTier() WorkspaceTier {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspaceTier == nil {
		repo.workspaceTier = &SqlWorkspaceTier{client: repo.client, timer: repo.timer}
	}

	return repo.workspaceTier
}

func (repo *sql) WorkspaceCredentials() WorkspaceCredentials {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspaceCredentials == nil {
		repo.workspaceCredentials = &SqlWorkspaceCredentials{client: repo.client, timer: repo.timer}
	}

	return repo.workspaceCredentials
}

func (repo *sql) Application() Application {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.client, timer: repo.timer}
	}

	return repo.application
}

func (repo *sql) Endpoint() Endpoint {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.client, timer: repo.timer}
	}

	return repo.endpoint
}

func (repo *sql) EndpointRule() EndpointRule {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpointRule == nil {
		repo.endpointRule = &SqlEndpointRule{client: repo.client, timer: repo.timer}
	}

	return repo.endpointRule
}
