package ds

import (
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
)

func New(logger logging.Logger, ds datastore.Datastore) Datastore {
	return NewSql(logger, ds)
}

type Datastore interface {
	Message() Message
	Request() Request
	Response() Response
}
