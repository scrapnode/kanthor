package dataplane

import (
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Controlplane interface {
	patterns.Connectable
}
