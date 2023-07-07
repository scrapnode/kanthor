package migration

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Migrator interface {
	Up() error
	Down() error
}

type Source interface {
	patterns.Connectable
	Migrator(source string) (Migrator, error)
}
