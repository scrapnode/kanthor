package repositories

import (
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
)

func New(logger logging.Logger, ds datastore.Datastore) Repositories {
	return NewSql(logger, ds)
}

type Repositories interface {
	Message() Message
	Request() Request
	Response() Response
}
