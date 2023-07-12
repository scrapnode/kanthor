package repositories

import (
	"context"
	xsql "database/sql"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlApplication) Create(ctx context.Context, app *entities.Application) (*entities.Application, error) {
	app.Id = utils.ID("app")
	app.CreatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Create(app); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return app, nil
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	if tx := sql.client.Model(&app).Where("id = ?", id).First(&app); tx.Error != nil {
		return nil, fmt.Errorf("application.get: %w", tx.Error)
	}

	if app.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("application.get.deleted: deleted_at:%d", app.DeletedAt)
	}

	return &app, nil
}

func (sql *SqlApplication) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error) {
	req := structure.ListReqBuild(opts)
	res := structure.ListRes[entities.Application]{}

	var tx = sql.client.Model(&entities.Application{}).
		Scopes(NotDeleted(sql.timer, &entities.Application{})).
		Where("workspace_id = ?", wsId)

	tx = TxListQuery(tx, req)
	if req.Search != "" {
		tx = tx.Where("name like ?", req.Search+"%")
	}

	if tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("application.list: %w", tx.Error)
	}

	return structure.ListResBuild(&res), nil
}

func (sql *SqlApplication) Update(ctx context.Context, app *entities.Application) (*entities.Application, error) {
	app.UpdatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Model(app).Select("name", "updated_at").Updates(app); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}

	return app, nil
}

func (sql *SqlApplication) Delete(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Non-repeatable_reads
	tx := sql.client.Begin(&xsql.TxOptions{Isolation: xsql.LevelReadCommitted})

	if txn := tx.Model(&app).Where("id = ?", id).First(&app); txn.Error != nil {
		return nil, fmt.Errorf("application.delete.get: %w", txn.Error)
	}

	app.UpdatedAt = sql.timer.Now().UnixMilli()
	app.DeletedAt = sql.timer.Now().UnixMilli()

	if txn := tx.Model(app).Select("updated_at", "deleted_at").Updates(app); txn.Error != nil {
		return nil, fmt.Errorf("application.delete.update: %w", txn.Error)
	}

	if txn := tx.Commit(); txn.Error != nil {
		return nil, fmt.Errorf("application.delete: %w", tx.Error)
	}

	return &app, nil
}
