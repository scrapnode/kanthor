package metric

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewNoop(conf *Config, logger logging.Logger) (Metrics, error) {
	logger = logger.With("monitoring.metrics", "noop")
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
func (metric *noop) Count(ctx context.Context, name string, value int64)     {}
func (metric *noop) Observe(ctx context.Context, name string, value float64) {}
