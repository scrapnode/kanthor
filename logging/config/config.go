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

	return &conf.Logger, nil
}

type Wrapper struct {
	Logger Config `json:"logger" yaml:"logger" mapstructure:"logger"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Logger.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Pretty bool              `json:"pretty" yaml:"pretty" mapstructure:"pretty"`
	Level  string            `json:"level" yaml:"level" mapstructure:"level"`
	With   map[string]string `json:"with" yaml:"with"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequiredOneOf("CONFIG.LOGGER.LEVEL", conf.Level, []string{"debug", "info", "warn", "error", "fatal"}),
	)
}
