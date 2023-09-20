package datastore

import (
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Uri       string           `json:"uri" yaml:"uri" mapstructure:"uri"`
	Migration migration.Config `json:"migration" yaml:"migration"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(validator.DefaultConfig, validator.StringUri("datastore.conf.uri", conf.Uri))
	if err != nil {
		return err
	}

	if err := conf.Migration.Validate(); err != nil {
		return err
	}

	return nil
}
