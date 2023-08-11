package database

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(conf *Config, logger logging.Logger, timer timer.Timer) (Database, error) {
	return NewSQL(conf, logger, timer), nil
}

type Database interface {
	patterns.Connectable
	Client() any
	Migrator(source string) (migration.Migrator, error)
}
