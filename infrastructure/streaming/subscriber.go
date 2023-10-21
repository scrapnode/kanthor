package streaming

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewSubscriber(conf *Config, logger logging.Logger) (Subscriber, error) {
	return NewNatsSubscriber(conf, logger), nil
}

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, name, topic string, handler SubHandler) error
}

type SubHandler func(events map[string]*Event) map[string]error
