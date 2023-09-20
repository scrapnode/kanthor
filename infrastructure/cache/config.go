package cache

import "github.com/scrapnode/kanthor/pkg/validator"

type Config struct {
	Uri        string `json:"uri" yaml:"uri" mapstructure:"uri"`
	TimeToLive int    `json:"time_to_live" yaml:"timeToLive" mapstructure:"time_to_live"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("idempotency.conf.uri", conf.Uri),
		validator.NumberGreaterThanOrEqual("idempotency.conf.time_to_live", conf.TimeToLive, 1000),
	)
}
