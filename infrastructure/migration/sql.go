package migration

import (
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
)

func NewSql(runner *migrate.Migrate) Migrator {
	return &sql{runner: runner}
}

type sql struct {
	runner *migrate.Migrate
}

func (migration *sql) Version() (uint, bool) {
	version, dirty, _ := migration.runner.Version()
	return version, dirty
}

func (migration *sql) Steps(n int) error {
	err := migration.runner.Steps(n)
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	// next/previous version is not exist
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}
