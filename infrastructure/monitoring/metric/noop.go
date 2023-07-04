package metric

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewNoop() Meter {
	return &prometheusio{}
}

func NewNoopExporter(conf *Config, logger logging.Logger) patterns.Runnable {
	return NewNoopServer(conf, logger)
}

type noop struct {
}

func (metric *noop) Counter(name string, value int64, labels ...Label)     {}
func (metric *noop) Histogram(name string, value float64, labels ...Label) {}
