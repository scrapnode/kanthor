package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
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
	tx := sql.client.WithContext(ctx).Model(&ws).Preload("Tier").Where("id = ?", id).First(&ws)
	if tx.Error != nil {
		return nil, fmt.Errorf("workspace.get: %w", tx.Error)
	}

	if ws.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("workspace.get.deleted: deleted_at:%d", ws.DeletedAt)
	}

	return &ws, nil
}

func (sql *SqlWorkspace) ListOfAccountSub(ctx context.Context, sub string, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error) {
	ws := &entities.Workspace{}
	wsc := &entities.WorkspacePrivilege{}
	join := fmt.Sprintf("JOIN %s ON %s.id = %s.workspace_id", ws.TableName(), ws.TableName(), wsc.TableName())
	selects := fmt.Sprintf("%s.*", ws.TableName())

	tx := sql.client.
		WithContext(ctx).
		Model(wsc).
		Joins(join).
		Where(fmt.Sprintf("%s.account_sub = ?", wsc.TableName()), sub).
		Scopes(database.NotDeleted(sql.timer, ws)).
		Scopes(database.NotDeleted(sql.timer, wsc)).
		Select(selects)
	tx = database.TxListQuery(tx, structure.ListReqBuild(opts))

	res := &structure.ListRes[entities.Workspace]{Data: []entities.Workspace{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return res, nil
}

func (sql *SqlWorkspace) GetByAccountSub(ctx context.Context, id, sub string) (*entities.Workspace, error) {
	ws := &entities.Workspace{}
	wsc := &entities.WorkspacePrivilege{}
	join := fmt.Sprintf("JOIN %s ON %s.id = %s.workspace_id", ws.TableName(), ws.TableName(), wsc.TableName())
	selects := fmt.Sprintf("%s.*", ws.TableName())

	tx := sql.client.
		WithContext(ctx).
		Model(wsc).
		Joins(join).
		Where(fmt.Sprintf("%s.id = ?", ws.TableName()), id).
		Where(fmt.Sprintf("%s.account_sub = ?", wsc.TableName()), sub).
		Scopes(database.NotDeleted(sql.timer, ws)).
		Scopes(database.NotDeleted(sql.timer, wsc)).
		Select(selects)

	if tx = tx.First(&ws); tx.Error != nil {
		return nil, tx.Error
	}

	return ws, nil
}
