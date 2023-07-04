package configuration

type Config interface {
	Validate() error
}
