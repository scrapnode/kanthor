package repos

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
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

	mu          sync.RWMutex
	client      *gorm.DB
	application *SqlApplication
	endpoint    *SqlEndpoint
	message     *SqlMessage
	request     *SqlRequest
	response    *SqlResponse
	attempt     *SqlAttempt
}

func (repo *sql) Readiness() error {
	if err := repo.db.Readiness(); err != nil {
		return err
	}
	return nil
}

func (repo *sql) Liveness() error {
	if err := repo.db.Liveness(); err != nil {
		return err
	}
	return nil
}

func (repo *sql) Connect(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.Connect(ctx); err != nil {
		return err
	}

	repo.client = repo.db.Client().(*gorm.DB)
	repo.logger.Info("connected")
	return nil
}

func (repo *sql) Disconnect(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.logger.Info("disconnected")

	if err := repo.db.Disconnect(ctx); err != nil {
		return err
	}

	repo.client = nil

	return nil
}

func (repo *sql) Application() Application {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.client}
	}

	return repo.application
}

func (repo *sql) Endpoint() Endpoint {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.client}
	}

	return repo.endpoint
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

func (repo *sql) Attempt() Attempt {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.attempt == nil {
		repo.attempt = &SqlAttempt{client: repo.client}
	}

	return repo.attempt
}
