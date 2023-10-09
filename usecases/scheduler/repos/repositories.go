package repos

import (
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(conf *database.Config, logger logging.Logger, timer timer.Timer) Repositories {
	return NewSql(conf, logger, timer)
}

type Repositories interface {
	patterns.Connectable
	Application() Application
	Endpoint() Endpoint
}
