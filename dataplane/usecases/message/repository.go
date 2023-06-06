package message

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"gorm.io/gorm"
)

type Repository interface {
	patterns.Connectable
	Put(ctx context.Context, message *entities.Message) (*entities.Message, error)
}

func NewRepository(logger logging.Logger, db datastore.Datastore) Repository {
	logger = logger.With("component", "dataplane.usecases.message.repository")
	return &SqlRepository{logger: logger, db: db}
}

type SqlRepository struct {
	logger logging.Logger
	db     datastore.Datastore
	orm    *gorm.DB
}

func (repo *SqlRepository) Connect(ctx context.Context) error {
	if err := repo.db.Connect(ctx); err != nil {
		return err
	}

	repo.orm = repo.db.DB().(*gorm.DB)

	repo.logger.Info("connected")
	return nil
}

func (repo *SqlRepository) Disconnect(ctx context.Context) error {
	repo.logger.Info("disconnected")

	if err := repo.db.Disconnect(ctx); err != nil {
		return err
	}

	repo.orm = nil

	return nil
}

func (repo *SqlRepository) Put(ctx context.Context, message *entities.Message) (*entities.Message, error) {
	if tx := repo.orm.Create(message); tx.Error != nil {
		return nil, fmt.Errorf("usecases.message.repository: %w", tx.Error)
	}

	return message, nil
}
