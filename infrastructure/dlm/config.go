package dlm

import "github.com/scrapnode/kanthor/pkg/validator"

type Config struct {
	Uri     string `json:"uri" yaml:"uri" mapstructure:"uri"`
	Timeout int    `json:"timeout" yaml:"timeToLive" mapstructure:"timeout"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("CONFIG.INFRA.DLM.URI", conf.Uri),
		validator.NumberGreaterThanOrEqual("CONFIG.INFRA.DLM.TIMEOUT", conf.Timeout, 1000),
	)
}
