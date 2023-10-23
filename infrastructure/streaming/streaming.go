package streaming

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Stream, error) {
	return NewNats(conf, logger)
}

type Stream interface {
	patterns.Connectable
	Publisher(name string) Publisher
	Subscriber(name string) Subscriber
}

type Publisher interface {
	Name() string
	Pub(ctx context.Context, events map[string]*Event) map[string]error
}

type Subscriber interface {
	patterns.Connectable
	Name() string
	Sub(ctx context.Context, topic string, handler SubHandler) error
}

type SubHandler func(events map[string]*Event) map[string]error
