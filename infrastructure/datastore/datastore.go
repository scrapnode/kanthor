package datastore

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) Datastore {
	return NewSQL(conf, logger)
}

type Datastore interface {
	patterns.Connectable
	Client() any
	Migrator(source string) (migration.Migrator, error)
}
