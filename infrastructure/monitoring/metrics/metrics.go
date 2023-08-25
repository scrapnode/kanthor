package metrics

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics/exporter"
)

func New(conf *Config, logger logging.Logger) (Metrics, error) {
	if conf.Engine == EngineNoop {
		return NewNoop(conf, logger)
	}
	if conf.Engine == EnginePrometheus {
		return NewPrometheus(conf, logger)
	}

	return nil, fmt.Errorf("authenticator: unknown engine")
}

type Metrics interface {
	Count(name string, value int64)
	Observe(name string, value float64)
	Exporter() exporter.Exporter
}
