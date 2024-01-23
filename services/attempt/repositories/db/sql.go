package db

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, db database.Database) Database {
	logger = logger.With("repositories", "db.sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     database.Database

	endpoint *SqlEndpoint

	mu sync.Mutex
}

func (repo *sql) Transaction(ctx context.Context, handler func(txctx context.Context) (interface{}, error)) (res interface{}, err error) {
	err = repo.db.Client().(*gorm.DB).Transaction(func(tx *gorm.DB) error {
		res, err = handler(context.WithValue(ctx, database.CtxTransaction, tx))
		return err
	})
	return
}

func (repo *sql) Endpoint() Endpoint {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.endpoint
}
