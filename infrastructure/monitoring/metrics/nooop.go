package metrics

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics/exporter"
)

func NewNoop(conf *Config, logger logging.Logger) (Metrics, error) {
	logger = logger.With("monitoring.metrics", "noop")
	return &noop{conf: conf, logger: logger}, nil
}

type noop struct {
	conf   *Config
	logger logging.Logger
}

func (metrics *noop) Count(name string, value int64)     {}
func (metrics *noop) Observe(name string, value float64) {}
func (metrics *noop) Exporter() exporter.Exporter {
	return exporter.NewNoop(
		&metrics.conf.Exporter,
		metrics.logger.With("metrics.exporter", "http"),
		nil,
	)
}
