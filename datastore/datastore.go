package datastore

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/datastore/config"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
)

func New(provider configuration.Provider) (Datastore, error) {
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

type Datastore interface {
	patterns.Connectable
	Client() any
}
