package repositories

import (
	"sync"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/scheduler/repositories/db"
)

func NewSql(logger logging.Logger, dbclient database.Database) Repositories {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db.NewSql(logger, dbclient)}
}

type sql struct {
	logger logging.Logger
	db     db.Database
	mu     sync.Mutex
}

func (repo *sql) Database() db.Database {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.db
}
