package authenticator

import (
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Engine  string   `json:"engine" yaml:"engine" mapstructure:"engine"`
	Ask     *Ask     `json:"ask" yaml:"ask" mapstructure:"ask"`
	Forward *Forward `json:"forward" yaml:"forward" mapstructure:"forward"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("AUTHENTICATOR.SCHEME", conf.Engine, []string{EngineAsk, EngineForward}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		return conf.Ask.Validate()
	}

	if conf.Engine == EngineForward {
		return conf.Forward.Validate()
	}

	return nil
}

type Ask struct {
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
}

func (conf *Ask) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("UTHENTICATOR.ASK.ACCESS_KEY", conf.AccessKey),
		validator.StringRequired("UTHENTICATOR.ASK.SECRET_KEY", conf.SecretKey),
	)
}

type Forward struct {
	Uri            string                `json:"uri" yaml:"uri" mapstructure:"uri"`
	Headers        []string              `json:"headers" yaml:"headers" mapstructure:"headers"`
	Sender         sender.Config         `json:"sender" yaml:"sender" mapstructure:"sender"`
	CircuitBreaker circuitbreaker.Config `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
}

func (conf *Forward) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("UTHENTICATOR.FORWARD.URI", conf.Uri),
	)
	if err != nil {
		return err
	}

	if err := conf.Sender.Validate(); err != nil {
		return err
	}

	if err := conf.CircuitBreaker.Validate(); err != nil {
		return err
	}

	return nil
}
