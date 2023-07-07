package patterns

type Migrate interface {
	Up() error
	Down() error
}

type Migrator interface {
	Connectable
	Migrator(source string) (Migrate, error)
}
