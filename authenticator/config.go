package authenticator

import (
	"errors"

	"github.com/scrapnode/kanthor/pkg/validator"
)

var (
	EngineAsk = "ask"
)

type Config struct {
	Engine string     `json:"engine" yaml:"engine" mapstructure:"engine"`
	Ask    *AskConfig `json:"ask" yaml:"ask" mapstructure:"ask"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("authenticator.config.engine", conf.Engine, []string{EngineAsk}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		if conf.Ask == nil {
			return errors.New("authenticator.config.ask: null value")
		}
		if err := conf.Ask.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type AskConfig struct {
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
}

func (conf *AskConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("authenticator.conf.ask.access_key", conf.AccessKey),
		validator.StringRequired("authenticator.conf.ask.secret_key", conf.SecretKey),
	)
}
