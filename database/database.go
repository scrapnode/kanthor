package database

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database/config"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
)

func New(provider configuration.Provider) (Database, error) {
	conf, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	return NewSQL(conf, logger)
}

type Database interface {
	patterns.Connectable
	Client() any
}
