package database

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Database, error) {
	return NewSQL(conf, logger)
}

type Database interface {
	patterns.Connectable
	DB() any
}
