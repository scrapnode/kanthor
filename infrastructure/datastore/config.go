package datastore

import (
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/migration"
)

type Config struct {
	Uri       string           `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
	Migration migration.Config `json:"migration" yaml:"migration" validate:"required"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if err := conf.Migration.Validate(); err != nil {
		return err
	}

	return nil
}
