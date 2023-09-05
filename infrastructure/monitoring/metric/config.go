package metric

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	EngineNoop = "noop"
	EngineOtel = "otel"
)

type Config struct {
	Engine string      `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=noop otel"`
	Otel   *OtelConfig `json:"otel" yaml:"otel" mapstructure:"otel" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
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
	Endpoint string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint" validate:"required,hostname_port"`
	Service  string            `json:"service" yaml:"service" mapstructure:"service" validate:"required"`
	Interval int               `json:"interval" yaml:"interval" mapstructure:"interval" validate:"required,gte=15000"`
	Labels   map[string]string `json:"labels" yaml:"labels" mapstructure:"labels" validate:"-"`
}

func (conf *OtelConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
