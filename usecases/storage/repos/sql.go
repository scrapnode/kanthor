package repos

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
	"sync"
)

func NewSql(conf *datastore.Config, logger logging.Logger, timer timer.Timer) Repositories {
	db := datastore.NewSQL(conf, logger, timer)

	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     datastore.Datastore

	mu       sync.RWMutex
	client   *gorm.DB
	message  *SqlMessage
	request  *SqlRequest
	response *SqlResponse
}

func (repo *sql) Connect(ctx context.Context) error {
	if err := repo.db.Connect(ctx); err != nil {
		return err
	}

	repo.client = repo.db.Client().(*gorm.DB)
	repo.logger.Info("connected")
	return nil
}

func (repo *sql) Disconnect(ctx context.Context) error {
	repo.logger.Info("disconnected")

	if err := repo.db.Disconnect(ctx); err != nil {
		return err
	}

	repo.client = nil

	return nil
}

func (repo *sql) Message() Message {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.message == nil {
		repo.message = &SqlMessage{client: repo.client}
	}

	return repo.message
}

func (repo *sql) Request() Request {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.request == nil {
		repo.request = &SqlRequest{client: repo.client}
	}

	return repo.request
}

func (repo *sql) Response() Response {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.response == nil {
		repo.response = &SqlResponse{client: repo.client}
	}

	return repo.response
}
