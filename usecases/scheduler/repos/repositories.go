package repos

import (
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func New(logger logging.Logger, db database.Database) Repositories {
	return NewSql(logger, db)
}

type Repositories interface {
	Application() Application
	Endpoint() Endpoint
}
