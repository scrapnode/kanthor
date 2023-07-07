package config

import "github.com/go-playground/validator/v10"

type Migration struct {
	Tasks []MigrationTask `json:"tasks" yaml:"tasks" mapstructure:"tasks" validate:"required"`
}

type MigrationTask struct {
	Name   string `json:"name" yaml:"name" mapstructure:"name" validate:"required"`
	Uri    string `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
	Source string `json:"source" yaml:"source" mapstructure:"source" validate:"required,uri"`
}

func (conf Migration) Validate() error {
	return validator.New().Struct(conf)
}
