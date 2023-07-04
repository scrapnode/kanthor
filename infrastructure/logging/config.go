package logging

import "github.com/go-playground/validator/v10"

type Config struct {
	Pretty bool              `json:"pretty" mapstructure:"pretty" validate:"boolean"`
	Level  string            `json:"level" mapstructure:"level" validate:"oneof=debug info warn error fatal"`
	With   map[string]string `json:"with" mapstructure:"with"`
}

func (conf Config) Validate() error {
	return validator.New().Struct(conf)
}
