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

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	var app entities.Application
	tx := transaction.WithContext(ctx).Model(&app).Where("id = ?", id).First(&app)
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("application.get: %w", err)
	}

	return &app, nil
}
