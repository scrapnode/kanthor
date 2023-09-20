package metric

import (
	"errors"

	"github.com/scrapnode/kanthor/pkg/validator"
)

var (
	EngineNoop = "noop"
	EngineOtel = "otel"
)

type Config struct {
	Engine string      `json:"engine" yaml:"engine" mapstructure:"engine"`
	Otel   *OtelConfig `json:"otel" yaml:"otel" mapstructure:"otel"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("monitoring.metrics.config.engine", conf.Engine, []string{EngineNoop, EngineOtel}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineOtel {
		if conf.Otel == nil {
			return errors.New("monitoring.metrics.config.otel: null value")
		}
		if err := conf.Otel.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type OtelConfig struct {
	Endpoint string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Service  string            `json:"service" yaml:"service" mapstructure:"service"`
	Interval int               `json:"interval" yaml:"interval" mapstructure:"interval"`
	Labels   map[string]string `json:"labels" yaml:"labels" mapstructure:"labels"`
}

func (conf *OtelConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringHostPort("monitoring.metrics.config.otel.endpoint", conf.Endpoint),
		validator.StringRequired("monitoring.metrics.config.otel.service", conf.Endpoint),
		validator.NumberGreaterThanOrEqual("monitoring.metrics.config.otel.interval", conf.Interval, 5000),
	)
}
