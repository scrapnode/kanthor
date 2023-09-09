package migration

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Migrator interface {
	// Version returns -1 mean there is no active version
	Version() (uint, bool)
	// Steps looks at the currently active migration version.
	// It will migrate up if n > 0, and down if n < 0.
	Steps(n int) error
}

type Source interface {
	patterns.Connectable
	Migrator(source string) (Migrator, error)
}
