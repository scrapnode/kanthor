package datastore

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/scrapnode/kanthor/datastore/config"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	postgresdevier "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSQL(conf *config.Config, logger logging.Logger) (Datastore, error) {
	logger = logger.With("datastore", "sql")
	return &sql{conf: conf, logger: logger}, nil
}

type sql struct {
	conf   *config.Config
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

	db, err := ds.client.DB()
	if err != nil {
		return err
	}

	// each postgres connection has their backend
	// the longer connection is alive, the more memory they consume
	db.SetConnMaxLifetime(time.Second * 300)
	db.SetConnMaxIdleTime(time.Second * 60)
	db.SetMaxIdleConns(1)

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
