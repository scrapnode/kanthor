package repositories

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/recovery/repositories/db"
	"github.com/scrapnode/kanthor/services/recovery/repositories/ds"
)

func New(logger logging.Logger, dbclient database.Database, dsclient datastore.Datastore) Repositories {
	return NewSql(logger, dbclient, dsclient)
}

type Repositories interface {
	Database() db.Database
	Datastore() ds.Datastore
}
