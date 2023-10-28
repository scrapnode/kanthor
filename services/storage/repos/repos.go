package repos

import (
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func New(logger logging.Logger, ds datastore.Datastore) Repositories {
	return NewSql(logger, ds)
}

type Repositories interface {
	Message() Message
	Request() Request
	Response() Response
}
