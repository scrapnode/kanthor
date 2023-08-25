package gateway

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

const EngineHttpx = "httpx"

type Config struct {
	Engine string       `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=httpx"`
	Httpx  *HttpxConfig `json:"httpx" yaml:"httpx" mapstructure:"httpx" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if conf.Engine == EngineHttpx {
		if conf.Httpx == nil {
			return errors.New("gateway.httpx: null value")
		}
		if err := conf.Httpx.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type HttpxConfig struct {
	Addr string `json:"addr" yaml:"addr" validate:"required"`
}

func (conf *HttpxConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
