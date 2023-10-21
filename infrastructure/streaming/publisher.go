package streaming

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewPublisher(conf *Config, logger logging.Logger) (Publisher, error) {
	return NewNatsPublisher(conf, logger), nil
}

type Publisher interface {
	patterns.Connectable
	Pub(ctx context.Context, events map[string]*Event) map[string]error
}
