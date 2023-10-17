package dlm

import "github.com/scrapnode/kanthor/pkg/validator"

type Config struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("dlm.conf.uri", conf.Uri),
	)
}
