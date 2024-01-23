package ds

import (
	"sync"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/logging"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, ds datastore.Datastore) Datastore {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, ds: ds}
}

type sql struct {
	logger logging.Logger
	ds     datastore.Datastore

	message  *SqlMessage
	request  *SqlRequest
	response *SqlResponse
	attempt  *SqlAttempt

	mu sync.Mutex
}

func (repo *sql) Message() Message {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.message == nil {
		repo.message = &SqlMessage{client: repo.ds.Client().(*gorm.DB)}
	}

	return repo.message
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

func (repo *sql) Attempt() Attempt {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.attempt == nil {
		repo.attempt = &SqlAttempt{client: repo.ds.Client().(*gorm.DB)}
	}

	return repo.attempt
}
