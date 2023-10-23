package logging

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Pretty bool              `json:"pretty" yaml:"pretty" mapstructure:"pretty"`
	Level  string            `json:"level" yaml:"level" mapstructure:"level"`
	With   map[string]string `json:"with" yaml:"with"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("logger.level", conf.Level, []string{"debug", "info", "warn", "error", "fatal"}),
	)
}
