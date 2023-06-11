package operators

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/migration/config"
)

func NewSql(conf *config.Config, logger logging.Logger) Operator {
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

func (operator *sql) Connect(ctx context.Context) error {
	if err := operator.db.Connect(ctx); err != nil {
		return err
	}

	m, err := operator.db.Migrator(operator.conf.Migration.Source)
	if err != nil {
		return err
	}
	operator.migrate = m

	operator.logger.Info("connected")
	return nil
}

func (operator *sql) Disconnect(ctx context.Context) error {
	operator.logger.Info("disconnected")

	if err := operator.db.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (operator *sql) Up() error {
	err := operator.migrate.Up()
	// convert error
	if errors.Is(err, migrate.ErrNoChange) {
		return ErrNoChange
	}

	return err
}

func (operator *sql) Down() error {
	return operator.migrate.Down()
}
