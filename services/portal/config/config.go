package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/gateway"
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

	return &conf.Portal, nil
}

type Wrapper struct {
	Portal Config `json:"portal" yaml:"portal" mapstructure:"portal"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Portal.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}

func (conf *Config) Validate() error {
	if err := conf.Gateway.Validate("CONFIG.PORTAL"); err != nil {
		return err
	}
	return nil
}
