package metric

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Metrics, error) {
	if conf.Engine == EngineNoop {
		return NewNoop(conf, logger)
	}
	if conf.Engine == EngineOtel {
		return NewOtel(conf, logger)
	}

	return nil, fmt.Errorf("authenticator: unknown engine")
}

type Metrics interface {
	patterns.Connectable
	Count(ctx context.Context, name string, value int64)
	Observe(ctx context.Context, name string, value float64)
}
