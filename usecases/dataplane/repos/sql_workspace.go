package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	var ws entities.Workspace
	tx := transaction.WithContext(ctx).Model(&ws).Preload("Tier").Where("id = ?", id).First(&ws)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ws, nil
}
