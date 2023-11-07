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

func (conf *Config) Validate(prefix string) error {
	if prefix != "" {
		prefix += "."
	}

	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf(prefix+"AUTHENTICATOR.ENGINE", conf.Engine, []string{EngineAsk}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		if conf.Ask == nil {
			return errors.New(prefix + "AUTHENTICATOR.ASK: nil value")
		}
		if err := conf.Ask.Validate(prefix); err != nil {
			return err
		}
	}

	return nil
}

type AskConfig struct {
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
}

func (conf *AskConfig) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("AUTHENTICATOR.ASK.ACCESS_KEY", conf.AccessKey),
		validator.StringRequired("AUTHENTICATOR.ASK.SECRET_KEY", conf.SecretKey),
	)
}
