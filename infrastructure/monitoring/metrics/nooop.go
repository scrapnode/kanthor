package metrics

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/http"
)

func NewNoop(conf *Config, logger logging.Logger) (Metrics, error) {
	logger = logger.With("monitoring.metrics", "noop")
	return &noop{}, nil
}

type noop struct {
}

func (metrics *noop) Count(name string, value int64)     {}
func (metrics *noop) Observe(name string, value float64) {}
func (metrics *noop) Handler() http.Handler              { return nil }
