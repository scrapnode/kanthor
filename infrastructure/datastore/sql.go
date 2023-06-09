package datastore

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"sync"
)

func NewSQL(conf *Config, logger logging.Logger) Datastore {
	return &sql{conf: conf, logger: logger.With("component", "datastore")}
}

type sql struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *gorm.DB
}

func (db *sql) Connect(ctx context.Context) error {
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

	db.logger.Info("connected")
	return err
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
