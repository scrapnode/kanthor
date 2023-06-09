package patterns

type Migrate interface {
	Up() error
	Down() error
}
