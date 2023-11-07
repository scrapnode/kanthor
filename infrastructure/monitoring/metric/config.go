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
		validator.StringOneOf("CONFIG.INFRA.MONITORING.METRIC.ENGINE", conf.Engine, []string{EngineNoop, EngineOtel}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineOtel {
		if conf.Otel == nil {
			return errors.New("CONFIG.INFRA.MONITORING.METRIC.OTEL: nil value")
		}
		if err := conf.Otel.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type OtelConfig struct {
	Endpoint string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Interval int               `json:"interval" yaml:"interval" mapstructure:"interval"`
	Labels   map[string]string `json:"labels" yaml:"labels" mapstructure:"labels"`
}

func (conf *OtelConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringHostPort("CONFIG.INFRA.MONITORING.METRIC.OTEL.ENDPOINT", conf.Endpoint),
		validator.NumberGreaterThanOrEqual("CONFIG.INFRA.MONITORING.METRIC.OTEL.INTERVAL", conf.Interval, 5000),
	)
}
