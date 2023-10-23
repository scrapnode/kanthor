package migration

type Migrator interface {
	// Version returns -1 mean there is no active version
	Version() (uint, bool)
	// Steps looks at the currently active migration version.
	// It will migrate up if n > 0, and down if n < 0.
	Steps(n int) error
}
