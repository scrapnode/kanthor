package repositories

import (
	"sync"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/storage/repositories/ds"
)

func NewSql(logger logging.Logger, dsclient datastore.Datastore) Repositories {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, ds: ds.NewSql(logger, dsclient)}
}

type sql struct {
	logger logging.Logger
	ds     ds.Datastore
	mu     sync.Mutex
}

func (repo *sql) Datastore() ds.Datastore {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.ds
}
