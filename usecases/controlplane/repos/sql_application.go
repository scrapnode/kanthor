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

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlApplication) Create(ctx context.Context, entity *entities.Application) (*entities.Application, error) {
	entity.GenId()
	entity.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entity); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return entity, nil
}

func (sql *SqlApplication) BulkCreate(ctx context.Context, entities []entities.Application) ([]string, error) {
	ids := []string{}
	for i, entity := range entities {
		entity.GenId()
		entity.SetAT(sql.timer.Now())

		ids = append(ids, entity.Id)
		entities[i] = entity
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entities); tx.Error != nil {
		return nil, fmt.Errorf("application.bulk_create: %w", tx.Error)
	}
	return ids, nil
}

func (sql *SqlApplication) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error) {
	app := &entities.Application{}
	tx := sql.client.
		WithContext(ctx).
		Model(app).
		Scopes(database.NotDeleted(sql.timer, app)).
		Where("workspace_id = ?", wsId)
	tx = database.SqlToListQuery(tx, structure.ListReqBuild(opts))

	res := &structure.ListRes[entities.Application]{Data: []entities.Application{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return res, nil
}

func (sql *SqlApplication) Get(ctx context.Context, wsId, id string) (*entities.Application, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	var app entities.Application
	tx := transaction.WithContext(ctx).Model(&app).
		Where("workspace_id = ?", wsId).
		Where("id = ?", id).
		First(&app)
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("application.get: %w", err)
	}

	if app.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("application.get.deleted: deleted_at:%d", app.DeletedAt)
	}

	return &app, nil
}
