package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
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
	Gateway       gateway.Config         `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Authenticator []authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator"`
}

func (conf *Config) Validate() error {
	if err := conf.Gateway.Validate("CONFIG.SDK"); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("CONFIG.SDK.AUTHENTICATOR", conf.Authenticator),
		validator.Slice(conf.Authenticator, func(i int, item *authenticator.Config) error {
			return item.Validate()
		}),
	)
}
