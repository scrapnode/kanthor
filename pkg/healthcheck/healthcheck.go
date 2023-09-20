package healthcheck

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Dest    string `json:"dest" yaml:"dest" mapstructure:"dest"`
	Timeout int    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	MaxTry  int    `json:"max_try" yaml:"max_try" mapstructure:"max_try"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("healthcheck.config.dest", conf.Dest),
		validator.NumberGreaterThanOrEqual("healthcheck.config.timeout", conf.Timeout, 0),
		validator.NumberGreaterThanOrEqual("healthcheck.config.max_try", conf.MaxTry, 0),
	)
}

func DefaultConfig(dest string) *Config {
	return &Config{Dest: fmt.Sprintf("kanthor.%s", dest), Timeout: 3000, MaxTry: 1}
}

type Server interface {
	Readiness(check func() error) error
	Liveness(check func() error) error
}

type Client interface {
	Readiness() error
	Liveness() error
}
