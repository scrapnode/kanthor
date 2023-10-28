package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Sdk Sdk `json:"sdk" yaml:"sdk" mapstructure:"sdk"`
}

func (conf *Config) Validate() error {
	if err := conf.Sdk.Validate(); err != nil {
		return err
	}
	return nil
}

type Sdk struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}

func (conf *Sdk) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("sdk.gateway: %v", err)
	}
	return nil
}
