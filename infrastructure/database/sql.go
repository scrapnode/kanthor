package database

import (
	"context"
	"errors"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	postgresdevier "gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewSQL(conf *Config, logger logging.Logger) Database {
	logger = logger.With("database", "sql")
	return &sql{conf: conf, logger: logger}
}

type sql struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *gorm.DB
}

func (db *sql) Readiness() error {
	if db.client == nil {
		return ErrNotConnected
	}

	var ok int
	tx := db.client.Raw("SELECT 1").Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotReady
	}

	return nil
}

func (db *sql) Liveness() error {
	if db.client == nil {
		return ErrNotConnected
	}

	var ok int
	tx := db.client.Raw("SELECT 1").Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotLive
	}

	return nil
}

func (db *sql) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		return ErrAlreadyConnected
	}

	dialector := postgresdevier.Open(db.conf.Uri)
	client, err := gorm.Open(dialector, &gorm.Config{
		Logger: NewSqlLogger(db.logger),
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
		conn, err := db.client.DB()
		if err != nil {
			return err
		}

		if err := conn.Close(); err != nil {
			return err
		}
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

	tableName := "kanthor_database_migration"
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

func SqlError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return err
}
