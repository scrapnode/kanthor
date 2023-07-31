package authenticator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	EngineAsk = "ask"
)

type Config struct {
	Engine string     `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=ask"`
	Ask    *AskConfig `json:"ask" yaml:"ask" mapstructure:"ask" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
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
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key" validate:"required"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key" validate:"required"`
}

func (conf *AskConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
