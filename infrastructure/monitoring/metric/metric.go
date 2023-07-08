package metric

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config) Meter {
	if !conf.Enable {
		return NewNoop(conf)
	}

	return NewPrometheus(conf)
}

func NewExporter(conf *Config, logger logging.Logger) patterns.Runnable {
	if !conf.Enable {
		return NewNoopExporter(conf, logger.With("monitoring.metrics.exporter", "noop"))
	}

	return NewPrometheusExporter(conf, logger.With("monitoring.metrics.exporter", "prometheus"))
}

type Meter interface {
	Count(name string, value int64, withLabels ...WithLabel)
	Histogram(name string, value float64, withLabels ...WithLabel)
}

type WithLabel func(labels map[string]string)

func Label(name, value string) WithLabel {
	return func(labels map[string]string) {
		labels[name] = value
	}
}
