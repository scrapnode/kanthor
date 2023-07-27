package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	var app entities.Application
	tx := transaction.WithContext(ctx).Model(&app).Where("id = ?", id).First(&app)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &app, nil
}
