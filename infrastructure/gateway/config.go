package gateway

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

const EngineHttpx = "httpx"

type Config struct {
	Addr string `json:"addr" yaml:"addr"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("gateway.config.addr", conf.Addr),
	)
}
