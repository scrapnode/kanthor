package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewSQL(conf *Config, logger logging.Logger) *SQL {
	return &SQL{conf: conf, logger: logger.With("component", "database")}
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
		dialector = sqlite.Open(uri.Host + uri.Path + uri.RawQuery)
	} else {
		dialector = postgres.Open(db.conf.Uri)
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
	driver, err := sqlite3.WithInstance(instance, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(source, "main", driver)
}
