package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf.Database, nil
}

type Wrapper struct {
	Database Config `json:"database" yaml:"database" mapstructure:"database"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Database.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Uri       string    `json:"uri" yaml:"uri" mapstructure:"uri"`
	Migration Migration `json:"migration" yaml:"migration"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("database.uri", conf.Uri),
		validator.StringStartsWith("database.uri", conf.Uri, "postgres"),
	)
	if err != nil {
		return err
	}

	if err := conf.Migration.Validate(); err != nil {
		return err
	}

	return nil
}

type Migration struct {
	Source string `json:"source" yaml:"source" mapstructure:"source"`
}

func (conf *Migration) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("database.migration.source", conf.Source),
	)
}