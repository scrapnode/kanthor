package healthcheck

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Dest      string      `json:"dest" yaml:"dest" mapstructure:"dest"`
	Readiness CheckConfig `json:"readiness" yaml:"readiness" mapstructure:"readiness"`
	Liveness  CheckConfig `json:"liveness" yaml:"liveness" mapstructure:"liveness"`
}

type CheckConfig struct {
	Timeout int `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	MaxTry  int `json:"max_try" yaml:"max_try" mapstructure:"max_try"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("healthcheck.config.dest", conf.Dest),
		validator.NumberGreaterThanOrEqual("healthcheck.config.readiness.timeout", conf.Readiness.Timeout, 0),
		validator.NumberGreaterThanOrEqual("healthcheck.config.readiness.max_try", conf.Readiness.MaxTry, 0),
		validator.NumberGreaterThanOrEqual("healthcheck.config.liveness.timeout", conf.Liveness.Timeout, 0),
		validator.NumberGreaterThanOrEqual("healthcheck.config.liveness.max_try", conf.Liveness.MaxTry, 0),
	)
}

func DefaultConfig(dest string) *Config {
	return &Config{
		Dest:      fmt.Sprintf("kanthor.%s", dest),
		Readiness: CheckConfig{Timeout: 10000, MaxTry: 3},
		Liveness:  CheckConfig{Timeout: 3000, MaxTry: 1},
	}
}

type Server interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Readiness(check func() error) error
	Liveness(check func() error) error
}

type Client interface {
	Readiness() error
	Liveness() error
}
