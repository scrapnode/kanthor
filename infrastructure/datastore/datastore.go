package datastore

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Datastore, error) {
	return NewSQL(conf, logger), nil
}

type Datastore interface {
	patterns.Connectable
	Client() any
	Migrator() (migration.Migrator, error)
}
