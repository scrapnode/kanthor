package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Portal Portal `json:"portal" yaml:"portal" mapstructure:"portal"`
}

func (conf *Config) Validate() error {
	if err := conf.Portal.Validate(); err != nil {
		return err
	}
	return nil
}

type Portal struct {
	Gateway       gateway.Config       `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Authenticator authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator"`
}

func (conf *Portal) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("portal.gateway: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("portal.authenticator: %v", err)
	}
	return nil
}
