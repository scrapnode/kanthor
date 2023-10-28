package metric

import (
	"context"

	"github.com/scrapnode/kanthor/logging"
)

func NewNoop(conf *Config, logger logging.Logger) (Metric, error) {
	logger = logger.With("metric", "noop")
	return &noop{conf: conf, logger: logger}, nil
}

type noop struct {
	conf   *Config
	logger logging.Logger
}

func (metric *noop) Readiness() error {
	return nil
}

func (metric *noop) Liveness() error {
	return nil
}

func (metric *noop) Connect(ctx context.Context) error {
	return nil
}

func (metric *noop) Disconnect(ctx context.Context) error {
	return nil
}

func (metric *noop) Count(ctx context.Context, service, name string, value int64) {}

func (metric *noop) Observe(ctx context.Context, service, name string, value float64) {}
