package migration

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Source string `json:"source" yaml:"source" mapstructure:"source"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("migration.config.engine", conf.Source),
	)
}
