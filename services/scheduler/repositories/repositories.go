package repositories

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
)

func New(logger logging.Logger, db database.Database) Repositories {
	return NewSql(logger, db)
}

type Repositories interface {
	Application() Application
	Endpoint() Endpoint
}
