package repositories

import (
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/storage/repositories/ds"
)

func New(logger logging.Logger, dsclient datastore.Datastore) Repositories {
	return NewSql(logger, dsclient)
}

type Repositories interface {
	Datastore() ds.Datastore
}
