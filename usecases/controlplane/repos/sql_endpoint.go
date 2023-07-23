package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlEndpoint) Create(ctx context.Context, entity *entities.Endpoint) (*entities.Endpoint, error) {
	entity.GenId()
	entity.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entity); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.create: %w", tx.Error)
	}
	return entity, nil
}

func (sql *SqlEndpoint) BulkCreate(ctx context.Context, entities []entities.Endpoint) ([]string, error) {
	ids := []string{}
	for i, entity := range entities {
		entity.GenId()
		entity.SetAT(sql.timer.Now())

		ids = append(ids, entity.Id)
		entities[i] = entity
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entities); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.bulk_create: %w", tx.Error)
	}
	return ids, nil
}
