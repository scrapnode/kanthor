package repos

import (
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *datastore.Config, logger logging.Logger) Repositories {
	return NewSql(conf, logger)
}

type Repositories interface {
	patterns.Connectable
	Application() Application
	Endpoint() Endpoint
	Message() Message
	Request() Request
	Response() Response
	Attempt() Attempt
}
