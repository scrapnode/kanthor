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

	client      *gorm.DB
	application *SqlApplication
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

func (repo *sql) Application() Application {
	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.client, timer: repo.timer}
	}

	return repo.application
}
