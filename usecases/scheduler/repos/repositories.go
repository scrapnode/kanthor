package repos

import (
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *database.Config, logger logging.Logger) Repositories {
	return NewSql(conf, logger)
}

type Repositories interface {
	patterns.Connectable
	Application() Application
	Endpoint() Endpoint
}
