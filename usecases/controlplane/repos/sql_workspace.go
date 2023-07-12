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
	var ws entities.Workspace
	tx := sql.client.Model(&ws).Preload("Tier").Where("id = ?", id).First(&ws)
	if tx.Error != nil {
		return nil, fmt.Errorf("workspace.get: %w", tx.Error)
	}

	if ws.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("workspace.get.deleted: deleted_at:%d", ws.DeletedAt)
	}

	return &ws, nil
}

func (sql *SqlWorkspace) ListByIds(ctx context.Context, ids []string) ([]entities.Workspace, error) {
	if len(ids) == 0 {
		return []entities.Workspace{}, nil
	}

	var tx = sql.client.Model(&entities.Workspace{}).
		Scopes(database.NotDeleted(sql.timer, &entities.Workspace{})).
		Where("id IN ?", ids)

	var workspaces []entities.Workspace
	if tx.Find(&workspaces); tx.Error != nil {
		return nil, fmt.Errorf("workspace.list: %w", tx.Error)
	}

	return workspaces, nil
}
