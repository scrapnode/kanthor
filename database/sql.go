package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/scrapnode/kanthor/database/config"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	postgresdevier "gorm.io/driver/postgres"
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
	if db.status == patterns.StatusDisconnected {
		return nil
	}
	if db.status != patterns.StatusConnected {
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

	if db.status == patterns.StatusConnected {
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
