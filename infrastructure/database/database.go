package database

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) Database {
	return NewSQL(conf, logger)
}

type Database interface {
	patterns.Connectable
	Client() any
	Migrator(source string) (patterns.Migrate, error)
}
