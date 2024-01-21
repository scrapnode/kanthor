package repositories

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/scheduler/repositories/db"
)

func New(logger logging.Logger, dbclient database.Database) Repositories {
	return NewSql(logger, dbclient)
}

type Repositories interface {
	Database() db.Database
}
