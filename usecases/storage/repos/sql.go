package repos

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"gorm.io/gorm"
)

func NewSql(conf *datastore.Config, logger logging.Logger) Repositories {
	db := datastore.NewSQL(conf, logger)

	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     datastore.Datastore

	client   *gorm.DB
	message  *SqlMessage
	request  *SqlRequest
	response *SqlResponse

	mu     sync.Mutex
	status int
}

func (repo *sql) Readiness() error {
	if repo.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := repo.db.Readiness(); err != nil {
		return err
	}
	return nil
}

func (repo *sql) Liveness() error {
	if repo.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := repo.db.Liveness(); err != nil {
		return err
	}
	return nil
}

func (repo *sql) Connect(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	if err := repo.db.Connect(ctx); err != nil {
		return err
	}
	repo.client = repo.db.Client().(*gorm.DB)

	repo.status = patterns.StatusConnected
	repo.logger.Info("connected")
	return nil
}

func (repo *sql) Disconnect(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	repo.status = patterns.StatusDisconnected
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
