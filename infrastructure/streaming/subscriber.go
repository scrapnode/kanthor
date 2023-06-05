package streaming

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, handler Handler) error
}

type Handler func(event *Event) error
