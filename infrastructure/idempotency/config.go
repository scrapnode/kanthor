package idempotency

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Uri        string `json:"uri" yaml:"uri" mapstructure:"uri"`
	TimeToLive uint64 `json:"time_to_live" yaml:"timeToLive" mapstructure:"time_to_live"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("CONFIG.INFRA.IDEMPOTENCY.URI", conf.Uri),
		validator.NumberGreaterThanOrEqual("CONFIG.INFRA.IDEMPOTENCY.TIME_TO_LIVE", conf.TimeToLive, 1000),
	)
}
