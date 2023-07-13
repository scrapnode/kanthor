package authenticator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Engine          string                 `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=access_secret_key"`
	AccessSecretKey *AccessSecretKeyConfig `json:"access_secret_key" yaml:"access_secret_key" mapstructure:"access_secret_key" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if conf.Engine == "access_secret_key" {
		if conf.AccessSecretKey == nil {
			return errors.New("authenticator.config.access_secret_key: null value")
		}
		if err := conf.AccessSecretKey.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type AccessSecretKeyConfig struct {
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key" validate:"required"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key" validate:"required"`
}

func (conf *AccessSecretKeyConfig) Validate() error {
	return validator.New().Struct(conf)
}
