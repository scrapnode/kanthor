package exporter

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Exporter interface {
	patterns.Runnable
}
