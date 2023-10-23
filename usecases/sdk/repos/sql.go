package repos

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"gorm.io/gorm"
)

func NewSql(conf *database.Config, logger logging.Logger) Repositories {
	db := database.NewSQL(conf, logger)

	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     database.Database

	client               *gorm.DB
	workspace            *SqlWorkspace
	workspaceCredentials *SqlWorkspaceCredentials
	application          *SqlApplication
	endpoint             *SqlEndpoint
	endpointRule         *SqlEndpointRule

	mu     sync.Mutex
	status int
}

func (repo *sql) Readiness() error {
	if repo.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := repo.db.Readiness(); err != nil {
		return err
	}
	return nil
}

func (repo *sql) Liveness() error {
	if repo.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := repo.db.Liveness(); err != nil {
		return err
	}
	return nil
}

func (repo *sql) Connect(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	if err := repo.db.Connect(ctx); err != nil {
		return err
	}
	repo.client = repo.db.Client().(*gorm.DB)

	repo.status = patterns.StatusConnected
	repo.logger.Info("connected")
	return nil
}

func (repo *sql) Disconnect(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	repo.status = patterns.StatusDisconnected
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
		repo.workspace = &SqlWorkspace{client: repo.client}
	}

	return repo.workspace
}

func (repo *sql) WorkspaceCredentials() WorkspaceCredentials {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspaceCredentials == nil {
		repo.workspaceCredentials = &SqlWorkspaceCredentials{client: repo.client}
	}

	return repo.workspaceCredentials
}

func (repo *sql) Application() Application {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.client}
	}

	return repo.application
}

func (repo *sql) Endpoint() Endpoint {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.client}
	}

	return repo.endpoint
}

func (repo *sql) EndpointRule() EndpointRule {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpointRule == nil {
		repo.endpointRule = &SqlEndpointRule{client: repo.client}
	}

	return repo.endpointRule
}
