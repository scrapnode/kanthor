package gateway

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

const EngineHttpx = "httpx"

type Config struct {
	Addr    string `json:"addr" yaml:"addr"`
	Timeout int64  `json:"timeout" yaml:"timeout"`
}

func (conf *Config) Validate(prefix string) error {
	if prefix != "" {
		prefix += "."
	}
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired(prefix+"GATEWAY.ADDR", conf.Addr),
		validator.NumberGreaterThanOrEqual(prefix+"GATEWAY.TIMEOUT", conf.Timeout, 1000),
	)
}
