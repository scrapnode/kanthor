package authenticator

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Bearer  *Bearer  `json:"bearer" yaml:"bearer" mapstructure:"bearer"`
	Forward *Forward `json:"forward" yaml:"forward" mapstructure:"forward"`
}

func (conf *Config) Validate() error {
	if conf.Bearer == nil && conf.Forward == nil {
		return fmt.Errorf("CONFIG.INFRA.AUTH.NO_SCHEME")
	}

	if conf.Bearer != nil {
		if err := conf.Bearer.Validate(); err != nil {
			return err
		}
	}

	if conf.Forward != nil {
		if err := conf.Forward.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Bearer struct {
	Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`
}

func (conf *Bearer) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringLen("CONFIG.INFRA.AUTH.BEARER.SECRET", conf.Secret, 16, 64),
	)
}

type Forward struct {
	Endpoint       string   `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	RequestHeaders []string `json:"request_headers" yaml:"request_headers" mapstructure:"request_headers"`
}

func (conf *Forward) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("CONFIG.INFRA.AUTH.FORWARD.ENDPOINT", conf.Endpoint),
	)
}
