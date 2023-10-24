package datastore

import (
	"context"
	"errors"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	postgresdevier "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSQL(conf *Config, logger logging.Logger) Datastore {
	logger = logger.With("datastore", "sql")
	return &sql{conf: conf, logger: logger}
}

type sql struct {
	conf   *Config
	logger logging.Logger

	client *gorm.DB

	mu     sync.Mutex
	status int
}

func (ds *sql) Readiness() error {
	if ds.status == patterns.StatusDisconnected {
		return nil
	}
	if ds.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := ds.client.Raw("SELECT 1").Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotReady
	}

	return nil
}

func (ds *sql) Liveness() error {
	if ds.status == patterns.StatusDisconnected {
		return nil
	}
	if ds.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := ds.client.Raw("SELECT 1").Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotLive
	}

	return nil
}

func (ds *sql) Connect(ctx context.Context) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	dialector := postgresdevier.Open(ds.conf.Uri)
	client, err := gorm.Open(dialector, &gorm.Config{
		// GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		// you can disable it during initialization if it is not required,
		// you will gain about 30%+ performance improvement after that
		SkipDefaultTransaction: true,
		Logger:                 NewSqlLogger(ds.logger),
	})
	if err != nil {
		return err
	}
	ds.client = client

	ds.status = patterns.StatusConnected
	ds.logger.Info("connected")
	return nil
}

func (ds *sql) Disconnect(ctx context.Context) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	ds.status = patterns.StatusDisconnected
	ds.logger.Info("disconnected")

	var returning error
	if conn, err := ds.client.DB(); err == nil {
		if err := conn.Close(); err != nil {
			returning = errors.Join(returning, err)
		}
	} else {
		returning = errors.Join(returning, err)
	}
	ds.client = nil

	return returning
}

func (ds *sql) Client() any {
	return ds.client
}

func (ds *sql) Migrator() (migration.Migrator, error) {
	instance, err := ds.client.DB()
	if err != nil {
		return nil, err
	}

	tableName := "kanthor_datastore_migration"
	var driver database.Driver

	conf := &postgres.Config{MigrationsTable: tableName}
	driver, err = postgres.WithInstance(instance, conf)
	if err != nil {
		return nil, err
	}

	runner, err := migrate.NewWithDatabaseInstance(ds.conf.Migration.Source, "", driver)
	if err != nil {
		return nil, err
	}

	return migration.NewSql(runner), nil
}
