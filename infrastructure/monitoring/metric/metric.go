package metric

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
)

func New(conf *Config, logger logging.Logger) (Metric, error) {
	if conf.Engine == EngineNoop {
		return NewNoop(conf, logger)
	}
	if conf.Engine == EngineOtel {
		return NewOtel(conf, logger)
	}

	return nil, fmt.Errorf("authenticator: unknown engine")
}

type Metric interface {
	patterns.Connectable
	Count(ctx context.Context, service, name string, value int64)
	Observe(ctx context.Context, service, name string, value float64)
}
