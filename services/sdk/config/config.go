package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/gateway"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Sdk, provider.Unmarshal(&conf)
}

type Wrapper struct {
	Sdk Config `json:"sdk" yaml:"sdk" mapstructure:"sdk"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Sdk.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}

func (conf *Config) Validate() error {
	if err := conf.Gateway.Validate("SDK.CONFIG."); err != nil {
		return err
	}

	return nil
}
