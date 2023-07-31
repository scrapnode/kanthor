package metric

import "github.com/go-playground/validator/v10"

const (
	EngineNoop = "noop"
)

type Config struct {
	Engine    string `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=noop"`
	Namespace string `json:"namespace" yaml:"namespace" mapstructure:"namespace" validate:"required"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
