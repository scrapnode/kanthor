package db

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
)

func New(logger logging.Logger, db database.Database) Database {
	return NewSql(logger, db)
}

type Database interface {
	Application() Application
}
