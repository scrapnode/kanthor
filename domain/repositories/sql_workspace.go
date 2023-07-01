package repositories

import (
	"context"
	xsql "database/sql"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspace) Create(ctx context.Context, ws *entities.Workspace) (*entities.Workspace, error) {
	ws.Id = utils.ID("ws")
	ws.CreatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Preload("Tier").Create(ws); tx.Error != nil {
		return nil, fmt.Errorf("workspace.create: %w", tx.Error)
	}
	return ws, nil
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

func (sql *SqlWorkspace) List(ctx context.Context, name string) ([]entities.Workspace, error) {
	var workspaces []entities.Workspace
	var tx = sql.client.Model(&entities.Workspace{}).
		Scopes(NotDeleted(sql.timer, &entities.Workspace{}))

	if name != "" {
		tx = tx.Where("name LIKE ?", name+"%")
	}

	if tx.Find(&workspaces); tx.Error != nil {
		return nil, fmt.Errorf("workspace"+
			".list: %w", tx.Error)
	}

	return workspaces, nil
}

func (sql *SqlWorkspace) Update(ctx context.Context, ws *entities.Workspace) (*entities.Workspace, error) {
	ws.UpdatedAt = sql.timer.Now().UnixMilli()

	tx := sql.client.Model(ws).
		Select("name", "updated_at").
		Updates(ws)
	if tx.Error != nil {
		return nil, fmt.Errorf("workspace.create: %w", tx.Error)
	}

	return ws, nil
}

func (sql *SqlWorkspace) Delete(ctx context.Context, id string) (*entities.Workspace, error) {
	var ws entities.Workspace
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Non-repreatable_reads
	tx := sql.client.Begin(&xsql.TxOptions{Isolation: xsql.LevelReadCommitted})

	if txn := tx.Model(&ws).Where("id = ?", id).First(&ws); txn.Error != nil {
		return nil, fmt.Errorf("workspace.delete.get: %w", txn.Error)
	}

	ws.UpdatedAt = sql.timer.Now().UnixMilli()
	ws.DeletedAt = sql.timer.Now().UnixMilli()

	if txn := tx.Model(ws).Select("updated_at", "deleted_at").Updates(ws); txn.Error != nil {
		return nil, fmt.Errorf("workspace.delete.update: %w", txn.Error)
	}

	if txn := tx.Commit(); txn.Error != nil {
		return nil, fmt.Errorf("workspace.delete: %w", tx.Error)
	}

	return &ws, nil
}
