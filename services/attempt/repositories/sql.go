package repositories

import (
	"sync"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services/attempt/repositories/db"
	"github.com/scrapnode/kanthor/services/attempt/repositories/ds"
)

func NewSql(logger logging.Logger, timer timer.Timer, dbclient database.Database, dsclient datastore.Datastore) Repositories {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db.NewSql(logger, dbclient), ds: ds.NewSql(logger, dsclient)}
}

type sql struct {
	logger logging.Logger

	db db.Database
	ds ds.Datastore
	mu sync.Mutex
}

func (repo *sql) Database() db.Database {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.db
}

func (repo *sql) Datastore() ds.Datastore {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.ds
}
