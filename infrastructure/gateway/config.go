package gateway

import (
	"errors"

	"github.com/scrapnode/kanthor/pkg/validator"
)

const EngineHttpx = "httpx"

type Config struct {
	Engine string       `json:"engine" yaml:"engine" mapstructure:"engine"`
	Httpx  *HttpxConfig `json:"httpx" yaml:"httpx" mapstructure:"httpx"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("gateway.config.engine", conf.Engine, []string{EngineHttpx}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineHttpx {
		if conf.Httpx == nil {
			return errors.New("gateway.config..httpx: null value")
		}
		if err := conf.Httpx.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type HttpxConfig struct {
	Addr string `json:"addr" yaml:"addr"`
}

func (conf *HttpxConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("gateway.config.addr", conf.Addr),
	)
}
