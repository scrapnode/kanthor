package gateway

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

const EngineHttpx = "httpx"

type Config struct {
	Addr    string `json:"addr" yaml:"addr"`
	Timeout int64  `json:"timeout" yaml:"timeout"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("gateway.config.addr", conf.Addr),
		validator.NumberGreaterThanOrEqual("gateway.config.timeout", conf.Timeout, 1000),
	)
}
