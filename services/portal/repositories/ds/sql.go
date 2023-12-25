package ds

import (
	"sync"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, ds datastore.Datastore, timer timer.Timer) Datastore {
	logger = logger.With("repositories", "datastore.sql")
	return &sql{logger: logger, ds: ds, timer: timer}
}

type sql struct {
	logger logging.Logger
	ds     datastore.Datastore
	timer  timer.Timer

	message *SqlMessage

	mu sync.Mutex
}

func (repo *sql) Message() Message {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.message == nil {
		repo.message = &SqlMessage{client: repo.ds.Client().(*gorm.DB), timer: repo.timer}
	}

	return repo.message
}
