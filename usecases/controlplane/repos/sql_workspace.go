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

func (sql *SqlWorkspace) Create(ctx context.Context, entity *entities.Workspace) (*entities.Workspace, error) {
	entity.GenId()
	entity.SetAT(sql.timer.Now())
	if entity.Tier != nil {
		entity.Tier.WorkspaceId = entity.Id
	}

	// if we use the entity to perform creating sql, gorm will use their logic to do the upsert
	// that leads to the error ON CONFLICT DO UPDATE requires inference specification or constraint name (SQLSTATE 42601)
	// so copy value and insert them separately is a workaround
	err := sql.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ws := &entities.Workspace{}
		ws.Id = entity.Id
		ws.CreatedAt = entity.CreatedAt
		ws.UpdatedAt = entity.UpdatedAt
		ws.DeletedAt = entity.DeletedAt
		ws.OwnerId = entity.OwnerId
		ws.Name = entity.Name

		if wstx := tx.Create(ws); wstx.Error != nil {
			return wstx.Error
		}

		tier := &entities.WorkspaceTier{WorkspaceId: entity.Tier.WorkspaceId, Name: entity.Tier.Name}
		if wsttx := tx.Create(tier); wsttx.Error != nil {
			return wsttx.Error
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return entity, nil
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

func (sql *SqlWorkspace) List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error) {
	ws := &entities.Workspace{}
	tx := sql.client.
		WithContext(ctx).
		Model(ws).
		Preload("Tier").
		Scopes(database.NotDeleted(sql.timer, ws))
	tx = database.TxListQuery(tx, structure.ListReqBuild(opts))

	res := &structure.ListRes[entities.Workspace]{Data: []entities.Workspace{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return res, nil
}
