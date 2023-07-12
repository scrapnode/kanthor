package dispatcher

import (
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Dispatcher interface {
	patterns.Connectable
	Forwarder() Forwarder
}
