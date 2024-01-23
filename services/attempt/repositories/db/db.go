package db

import (
	"context"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
)

func New(logger logging.Logger, db database.Database) Database {
	return NewSql(logger, db)
}

type Database interface {
	Transaction(ctx context.Context, handler func(txctx context.Context) (interface{}, error)) (res interface{}, err error)
	Endpoint() Endpoint
}
