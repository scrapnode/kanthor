package datastore

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) Datastore {
	return NewSQL(conf, logger)
}

type Datastore interface {
	patterns.Connectable
	DB() any
}
