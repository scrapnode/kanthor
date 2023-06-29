package migration

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services"
)

func New(conf *config.Config, logger logging.Logger, db database.Database) services.Service {
	logger.With("service", "migration")
	return &migration{conf: conf, logger: logger, db: db}
}

type migration struct {
	conf   *config.Config
	logger logging.Logger
	db     database.Database

	migrate patterns.Migrate
}

func (service *migration) Start(ctx context.Context) error {
	if err := service.db.Connect(ctx); err != nil {
		return err
	}

	m, err := service.db.Migrator(service.conf.Migration.Source)
	if err != nil {
		return err
	}
	service.migrate = m

	service.logger.Info("connected")
	return nil
}

func (service *migration) Stop(ctx context.Context) error {
	service.logger.Info("disconnected")

	if err := service.db.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (service *migration) Run(ctx context.Context) error {
	err := service.migrate.Up()

	// ignore no change error
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
