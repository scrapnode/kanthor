package message

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type Repository interface {
	Put(ctx context.Context, message *entities.Message) (*entities.Message, error)
}

func NewRepository(db database.Database) Repository {
	return &SqlRepository{db: db.DB().(*gorm.DB)}
}

type SqlRepository struct {
	db *gorm.DB
}

func (repo *SqlRepository) Put(ctx context.Context, message *entities.Message) (*entities.Message, error) {
	if tx := repo.db.Create(message); tx.Error != nil {
		return nil, fmt.Errorf("application.message.repository: %w", tx.Error)
	}

	return message, nil
}
