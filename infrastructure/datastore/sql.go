package datastore

import (
	"context"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/pkg/timer"
	postgresdevier "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSQL(conf *Config, logger logging.Logger, timer timer.Timer) Datastore {
	logger = logger.With("datastore", "sql")
	return &sql{conf: conf, logger: logger, timer: timer}
}

type sql struct {
	conf   *Config
	logger logging.Logger
	timer  timer.Timer

	mu     sync.Mutex
	client *gorm.DB
}

func (db *sql) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		return ErrAlreadyConnected
	}

	dialector := postgresdevier.Open(db.conf.Uri)
	client, err := gorm.Open(dialector, &gorm.Config{
		// GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		// you can disable it during initialization if it is not required,
		// you will gain about 30%+ performance improvement after that
		SkipDefaultTransaction: true,
		Logger:                 NewSqlLogger(db.logger),
		NowFunc: func() time.Time {
			return db.timer.Now()
		},
		// @TODO: considering options
		// CreateBatchSize: 100,
		// PrepareStmt: false,
	})
	if err != nil {
		return err
	}
	db.client = client

	db.logger.Info("connected")
	return nil
}

func (db *sql) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client == nil {
		return ErrNotConnected
	}

	conn, err := db.client.DB()
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	db.client = nil
	db.logger.Info("disconnected")
	return nil
}

func (db *sql) Client() any {
	return db.client
}

func (db *sql) Migrator() (migration.Migrator, error) {
	instance, err := db.client.DB()
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

	runner, err := migrate.NewWithDatabaseInstance(db.conf.Migration.Source, "", driver)
	if err != nil {
		return nil, err
	}

	return migration.NewSql(runner), nil
}
