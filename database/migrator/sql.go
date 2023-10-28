package migrator

import (
	dbsql "database/sql"
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/scrapnode/kanthor/database/config"
	"github.com/scrapnode/kanthor/project"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewSql(conf *config.Config) (Migrator, error) {
	db, err := dbsql.Open("postgres", conf.Uri)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: project.NameWithoutTier("database_migration")})
	if err != nil {
		return nil, err
	}

	runner, err := migrate.NewWithDatabaseInstance(conf.Migration.Source, "", driver)
	if err != nil {
		return nil, err
	}

	return &sql{runner: runner}, nil
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
