package migration

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
)

func NewSql(runner *migrate.Migrate) Migrator {
	return &sql{runner: runner}
}

type sql struct {
	runner *migrate.Migrate
}

func (migration *sql) Up() error {
	err := migration.runner.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (migration *sql) Down() error {
	return migration.runner.Down()
}
