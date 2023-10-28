package migrator

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/datastore/config"
)

func New(provider configuration.Provider) (Migrator, error) {
	conf, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	return NewSql(conf)
}

type Migrator interface {
	// Version returns -1 mean there is no active version
	Version() (uint, bool)
	// Steps looks at the currently active migration version.
	// It will migrate up if n > 0, and down if n < 0.
	Steps(n int) error
}
