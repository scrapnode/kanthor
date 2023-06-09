package migrators

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/migration/config"
)

func NewSql(conf *config.Config, logger logging.Logger) Migrator {
	db := database.New(&conf.Database, logger)

	logger = logger.With("component", "migration.sql")
	return &sql{conf: conf, logger: logger, db: db}
}

type sql struct {
	conf   *config.Config
	logger logging.Logger
	db     database.Database

	migrate patterns.Migrate
}

func (migrator *sql) Connect(ctx context.Context) error {
	if err := migrator.db.Connect(ctx); err != nil {
		return err
	}

	m, err := migrator.db.Migrator(migrator.conf.Migration.Source)
	if err != nil {
		return err
	}
	migrator.migrate = m

	migrator.logger.Info("connected")
	return nil
}

func (migrator *sql) Disconnect(ctx context.Context) error {
	migrator.logger.Info("disconnected")

	if err := migrator.db.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (migrator *sql) Up() error {
	err := migrator.migrate.Up()
	// convert error
	if errors.Is(err, migrate.ErrNoChange) {
		return ErrNoChange
	}

	return err
}

func (migrator *sql) Down() error {
	return migrator.migrate.Down()
}
