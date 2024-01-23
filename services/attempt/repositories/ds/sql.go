package ds

import (
	"sync"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, ds datastore.Datastore) Datastore {
	logger = logger.With("repositories", "datastore.sql")
	return &sql{logger: logger, ds: ds}
}

type sql struct {
	logger logging.Logger
	ds     datastore.Datastore

	request  *SqlRequest
	response *SqlResponse

	mu sync.Mutex
}

func (repo *sql) Request() Request {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.request == nil {
		repo.request = &SqlRequest{client: repo.ds.Client().(*gorm.DB)}
	}

	return repo.request
}

func (repo *sql) Response() Response {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.response == nil {
		repo.response = &SqlResponse{client: repo.ds.Client().(*gorm.DB)}
	}

	return repo.response
}
