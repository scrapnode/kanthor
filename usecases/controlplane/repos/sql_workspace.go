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

// SqlWorkspace
// SqlClientFromContext must be used in all single object action to reuse them in other place: Create, Get, GetOwned
// Start transaction in listing produces a risk of un-predictable locking row, so we should not use SqlClientFromContext
type SqlWorkspace struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspace) Create(ctx context.Context, entity *entities.Workspace) (*entities.Workspace, error) {
	now := sql.timer.Now()

	entity.GenId()
	entity.SetAT(now)

	entity.Tier.WorkspaceId = entity.Id
	entity.Tier.GenId()
	entity.Tier.SetAT(now)

	transaction := database.SqlClientFromContext(ctx, sql.client)
	// if we use the entity to perform creating sql, gorm will use their logic to do the upsert
	// that leads to the error ON CONFLICT DO UPDATE requires inference specification or constraint name (SQLSTATE 42601)
	// so copy value and insert them separately is a workaround
	// NOTE: sub-transaction is needed to use in other places
	err := transaction.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ws := &entities.Workspace{}
		ws.Id = entity.Id
		ws.ModifiedBy = entity.ModifiedBy
		ws.CreatedAt = entity.CreatedAt
		ws.UpdatedAt = entity.UpdatedAt
		ws.OwnerId = entity.OwnerId
		ws.Name = entity.Name

		if wstx := tx.Create(ws); wstx.Error != nil {
			return wstx.Error
		}

		tier := &entities.WorkspaceTier{WorkspaceId: entity.Tier.WorkspaceId, Name: entity.Tier.Name}
		tier.Id = entity.Tier.Id
		tier.ModifiedBy = entity.Tier.ModifiedBy
		tier.CreatedAt = entity.Tier.CreatedAt
		tier.UpdatedAt = entity.Tier.UpdatedAt

		if wsttx := tx.Create(tier); wsttx.Error != nil {
			return wsttx.Error
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("workspace.create: %w", err)
	}

	return entity, nil
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	var ws entities.Workspace
	tx := transaction.WithContext(ctx).Model(&ws).Preload("Tier").Where("id = ?", id).First(&ws)
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("workspace.get: %w", err)
	}

	return &ws, nil
}

func (sql *SqlWorkspace) GetOwned(ctx context.Context, owner string) (*entities.Workspace, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	var ws entities.Workspace
	tx := transaction.WithContext(ctx).Model(&ws).Preload("Tier").Where("owner_id = ?", owner).First(&ws)
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("workspace.get_owned: %w", err)
	}

	return &ws, nil
}

func (sql *SqlWorkspace) List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error) {
	req := structure.ListReqBuild(opts)

	ws := &entities.Workspace{}

	tx := sql.client.
		WithContext(ctx).
		Model(ws).
		Preload("Tier")
	tx = database.SqlToListQuery(tx, req, `"id"`)

	res := &structure.ListRes[entities.Workspace]{Data: []entities.Workspace{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("workspace.list: %w", tx.Error)
	}

	return res, nil
}
