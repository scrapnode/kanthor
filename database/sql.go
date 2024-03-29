package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/scrapnode/kanthor/database/config"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewSQL(conf *config.Config, logger logging.Logger) (Database, error) {
	logger = logger.With("database", "sql")
	return &sql{conf: conf, logger: logger}, nil
}

type sql struct {
	conf   *config.Config
	logger logging.Logger

	client *gorm.DB

	mu     sync.Mutex
	status int
}

func (db *sql) Readiness() error {
	if db.status == patterns.StatusDisconnected {
		return nil
	}
	if db.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := db.client.Raw(readinessQuery).Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotReady
	}

	return nil
}

func (db *sql) Liveness() error {
	if db.status == patterns.StatusDisconnected {
		return nil
	}
	if db.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := db.client.Raw(livenessQuery).Scan(&ok)
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

	if db.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	client, err := db.driver()
	if err != nil {
		return err
	}
	db.client = client

	isntance, err := db.client.DB()
	if err != nil {
		return err
	}
	// each postgres connection has their backend
	// the longer connection is alive, the more memory they consume
	isntance.SetConnMaxLifetime(time.Second * 300)
	isntance.SetConnMaxIdleTime(time.Second * 60)
	isntance.SetMaxIdleConns(1)
	isntance.SetMaxOpenConns(10)

	db.status = patterns.StatusConnected
	db.logger.Info("connected")
	return nil
}

func (db *sql) driver() (*gorm.DB, error) {
	dialector := postgresdriver.Open(db.conf.Uri)
	return gorm.Open(dialector, &gorm.Config{
		Logger: NewSqlLogger(db.logger),
	})
}

func (db *sql) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	db.status = patterns.StatusDisconnected
	db.logger.Info("disconnected")

	var returning error
	if conn, err := db.client.DB(); err == nil {
		if err := conn.Close(); err != nil {
			returning = errors.Join(returning, err)
		}
	} else {
		returning = errors.Join(returning, err)
	}
	db.client = nil

	return returning
}

func (db *sql) Client() any {
	return db.client
}

var (
	readinessQuery = "SELECT 1 as readiness"
	livenessQuery  = "SELECT 1 as liveness"
)
