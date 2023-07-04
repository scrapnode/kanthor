package metric

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		return NewNoopServer(conf, logger.With("monitoring.metrics.exporter", "noop"))
	}

	return NewHttpServer(conf, logger.With("monitoring.metrics.exporter", "prometheus"), promhttp.Handler())
}

type Meter interface {
	Counter(name string, value int64, labels ...Label)
	Histogram(name string, value float64, labels ...Label)
}

type Label func(labels map[string]string)

func UseLabel(name, value string) Label {
	return func(labels map[string]string) {
		labels[name] = value
	}
}
