package db

import (
	"sync"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/logging"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, db database.Database) Database {
	logger = logger.With("component", "repositories.sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     database.Database

	application *SqlApplication

	mu sync.Mutex
}

func (repo *sql) Application() Application {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.application
}
