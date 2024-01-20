package migrator

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/scrapnode/kanthor/datastore/config"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/project"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var TableMigration = "datastore_migration"

func NewSql(conf *config.Config) (Migrator, error) {
	scheme, d, err := driver(conf)
	if err != nil {
		return nil, err
	}

	runner, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("%s/%s", conf.Migration.Source, scheme), "", d)
	if err != nil {
		return nil, err
	}

	return &sql{runner: runner}, nil
}

func driver(conf *config.Config) (string, database.Driver, error) {
	scheme, err := utils.UrlScheme(conf.Uri)
	if err != nil {
		return "", nil, err
	}

	db, err := dbsql.Open("postgres", conf.Uri)
	if err != nil {
		return "", nil, err
	}
	d, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: project.Name(TableMigration)})
	return scheme, d, err
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
