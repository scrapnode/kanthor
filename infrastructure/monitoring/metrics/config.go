package metrics

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics/exporter"
)

var (
	EngineNoop       = "noop"
	EnginePrometheus = "prometheus"
)

type Config struct {
	Exporter   exporter.Config   `json:"exporter" yaml:"exporter" mapstructure:"exporter" validate:"required"`
	Engine     string            `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=noop prometheus"`
	Prometheus *PrometheusConfig `json:"prometheus" yaml:"prometheus" mapstructure:"prometheus" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if err := conf.Exporter.Validate(); err != nil {
		return err
	}

	if conf.Engine == EnginePrometheus {
		if conf.Prometheus == nil {
			return errors.New("monitoring.metrics.config.prometheus: null value")
		}
		if err := conf.Prometheus.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type PrometheusConfig struct {
	Labels map[string]string `json:"labels" yaml:"labels" mapstructure:"labels" validate:"-"`
}

func (conf *PrometheusConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
