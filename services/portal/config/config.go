package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/authenticator"
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
	Gateway       gateway.Config       `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Authenticator authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator"`
}

func (conf *Config) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("portal.gateway: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("portal.authenticator: %v", err)
	}
	return nil
}
