package repos

import (
	"context"
	"fmt"
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
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("workspace.get: %w", err)
	}

	if ws.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("workspace.get.deleted: deleted_at:%d", ws.DeletedAt)
	}

	return &ws, nil
}
