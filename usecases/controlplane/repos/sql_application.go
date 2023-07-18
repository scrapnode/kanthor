package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlApplication) Create(ctx context.Context, entity *entities.Application) (*entities.Application, error) {
	entity.GenId()
	entity.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entity); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return entity, nil
}

func (sql *SqlApplication) BulkCreate(ctx context.Context, entities []entities.Application) ([]string, error) {
	var ids []string
	for i, entity := range entities {
		entity.GenId()
		entity.SetAT(sql.timer.Now())

		ids = append(ids, entity.Id)
		entities[i] = entity
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entities); tx.Error != nil {
		return nil, fmt.Errorf("application.bulk_create: %w", tx.Error)
	}
	return ids, nil
}