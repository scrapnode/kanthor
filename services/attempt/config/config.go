package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Attempt, provider.Unmarshal(&conf)
}

type Wrapper struct {
	Attempt Config `json:"attempt" yaml:"attempt" mapstructure:"attempt"`
}

type Config struct {
	Trigger  Trigger  `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
	Endeavor Endeavor `json:"endeavor" yaml:"endeavor" mapstructure:"endeavor"`
}

func (conf *Config) Validate() error {
	if err := conf.Trigger.Validate(); err != nil {
		return fmt.Errorf("attempt.trigger: %v", err)
	}

	if err := conf.Endeavor.Validate(); err != nil {
		return fmt.Errorf("attempt.endeavor: %v", err)
	}

	return nil
}
