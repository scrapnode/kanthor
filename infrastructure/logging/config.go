package logging

import "github.com/go-playground/validator/v10"

type Config struct {
	Pretty bool              `json:"pretty" yaml:"pretty" mapstructure:"pretty" validate:"boolean"`
	Level  string            `json:"level" yaml:"level" mapstructure:"level" validate:"oneof=debug info warn error fatal"`
	With   map[string]string `json:"with" yaml:"with" mapstructure:"with"`
}

func (conf Config) Validate() error {
	return validator.New().Struct(conf)
}
