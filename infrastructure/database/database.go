package database

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Database, error) {
	return NewSQL(conf, logger), nil
}

type Database interface {
	patterns.Connectable
	Client() any
	Migrator() (migration.Migrator, error)
}
