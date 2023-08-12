package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	if tx := sql.client.Model(&app).Where("id = ?", id).First(&app); tx.Error != nil {
		return nil, fmt.Errorf("application.get: %w", tx.Error)
	}

	return &app, nil
}
