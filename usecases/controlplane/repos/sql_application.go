package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
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
	
	if tx := sql.client.WithContext(ctx).Create(entity); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return entity, nil
}
