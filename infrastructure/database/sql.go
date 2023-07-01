package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	postgresdevier "gorm.io/driver/postgres"
	sqlitedriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewSQL(conf *Config, logger logging.Logger) *SQL {
	logger = logger.With("component", "database.sql")
	return &SQL{conf: conf, logger: logger}
}

type SQL struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *gorm.DB
}

func (db *SQL) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		return ErrAlreadyConnected
	}

	uri, err := url.Parse(db.conf.Uri)
	if err != nil {
		return err
	}

	var dialector gorm.Dialector
	if strings.HasPrefix(uri.Scheme, "sqlite") {
		dialector = sqlitedriver.Open(uri.Host + uri.Path)
	} else {
		dialector = postgresdevier.Open(db.conf.Uri)
	}

	db.client, err = gorm.Open(dialector, &gorm.Config{Logger: NewSqlLogger(db.logger)})
	if err != nil {
		return err
	}

	db.logger.Info("connected")
	return nil
}

func (db *SQL) Disconnect(ctx context.Context) error {
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

func (db *SQL) Client() any {
	return db.client
}

func (db *SQL) Migrator(source string) (patterns.Migrate, error) {
	instance, err := db.client.DB()
	if err != nil {
		return nil, err
	}

	var driver database.Driver
	if db.client.Config.Dialector.Name() == "sqlite" {
		conf := &sqlite3.Config{
			MigrationsTable: "database_migration",
		}
		driver, err = sqlite3.WithInstance(instance, conf)
		if err != nil {
			return nil, err
		}
	} else {
		conf := &postgres.Config{
			MigrationsTable: "database_migration",
		}
		driver, err = postgres.WithInstance(instance, conf)
		if err != nil {
			return nil, err
		}
	}

	return migrate.NewWithDatabaseInstance(source, "", driver)
}
