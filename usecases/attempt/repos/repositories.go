package repos

import (
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func New(logger logging.Logger, ds datastore.Datastore) Repositories {
	return NewSql(logger, ds)
}

type Repositories interface {
	Application() Application
	Endpoint() Endpoint
	Message() Message
	Request() Request
	Response() Response
	Attempt() Attempt
}
