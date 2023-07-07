package config

import "github.com/go-playground/validator/v10"

type Migration struct {
	Source string `json:"source" yaml:"source" mapstructure:"source" validate:"required,uri"`
}

func (conf Migration) Validate() error {
	return validator.New().Struct(conf)
}
