package dlm

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Uri        string `json:"uri" yaml:"uri" mapstructure:"uri"`
	TimeToLive uint64 `json:"time_to_live" yaml:"time_to_live" mapstructure:"time_to_live"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("INFRASTRUCTURE.DLM.CONFIG.URI", conf.Uri),
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.DLM.CONFIG.TIME_TO_LIVE", conf.TimeToLive, 1000),
	)
}

type Option func(*Config)

func TimeToLive(ttl uint64) Option {
	return func(conf *Config) {
		conf.TimeToLive = ttl
	}
}
