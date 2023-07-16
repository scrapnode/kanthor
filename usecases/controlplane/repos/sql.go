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

	client       *gorm.DB
	workspace    *SqlWorkspace
	application  *SqlApplication
	endpoint     *SqlEndpoint
	endpointRule *SqlEndpointRule
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

	return nil
}

func (repo *sql) Workspace() Workspace {
	if repo.workspace == nil {
		repo.workspace = &SqlWorkspace{client: repo.client, timer: repo.timer}
	}

	return repo.workspace
}

func (repo *sql) Application() Application {
	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.client, timer: repo.timer}
	}

	return repo.application
}

func (repo *sql) Endpoint() Endpoint {
	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.client, timer: repo.timer}
	}

	return repo.endpoint
}

func (repo *sql) EndpointRule() EndpointRule {
	if repo.endpointRule == nil {
		repo.endpointRule = &SqlEndpointRule{client: repo.client, timer: repo.timer}
	}

	return repo.endpointRule
}
