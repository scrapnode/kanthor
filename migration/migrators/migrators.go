package migrators

import (
	"errors"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/migration/config"
)

func New(conf *config.Config, logger logging.Logger) Migrator {
	return NewSql(conf, logger)
}

type Migrator interface {
	patterns.Connectable
	Up() error
	Down() error
}

var ErrNoChange = errors.New("migration: no change")
