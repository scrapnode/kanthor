package authenticator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
)

var (
	EngineAsk    = "ask"
	EngineCipher = "cipher"
)

type Config struct {
	Engine string                        `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=ask cipher"`
	Ask    *AskConfig                    `json:"ask" yaml:"ask" mapstructure:"ask" validate:"-"`
	Cipher *cryptography.SymmetricConfig `json:"cipher" yaml:"cipher" mapstructure:"cipher" validate:"-"`
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

	if conf.Engine == EngineCipher {
		if conf.Cipher == nil {
			return errors.New("authenticator.config.cipher: null value")
		}
		if err := conf.Cipher.Validate(); err != nil {
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
