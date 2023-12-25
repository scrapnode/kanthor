package repositories

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services/attempt/repositories/db"
	"github.com/scrapnode/kanthor/services/attempt/repositories/ds"
)

func New(logger logging.Logger, timer timer.Timer, dbclient database.Database, dsclient datastore.Datastore) Repositories {
	return NewSql(logger, timer, dbclient, dsclient)
}

type Repositories interface {
	Database() db.Database
	Datastore() ds.Datastore
}
