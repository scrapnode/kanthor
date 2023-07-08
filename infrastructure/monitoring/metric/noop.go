package metric

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewNoop(conf *Config) Meter {
	return &noop{}
}

func NewNoopExporter(conf *Config, logger logging.Logger) patterns.Runnable {
	return NewNoopServer(conf, logger)
}

type noop struct {
}

func (metric *noop) Count(name string, value int64, withLabels ...WithLabel)       {}
func (metric *noop) Histogram(name string, value float64, withLabels ...WithLabel) {}
