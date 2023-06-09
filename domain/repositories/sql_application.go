package repositories

import (
	"context"
	xsql "database/sql"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/timer"
	"github.com/scrapnode/kanthor/infrastructure/utils"
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
		return nil, fmt.Errorf("repositories.sql.application.create: %w", tx.Error)
	}
	return app, nil
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	if tx := sql.client.Model(&app).Where("id = ?", id).First(&app); tx.Error != nil {
		return nil, fmt.Errorf("repositories.sql.application.get: %w", tx.Error)
	}

	return &app, nil
}

func (sql *SqlApplication) List(ctx context.Context, wsId, name string) ([]entities.Application, error) {
	var apps []entities.Application
	var tx = sql.client.Model(&entities.Application{}).Where("workspace_id = ?", wsId)

	if name != "" {
		tx = tx.Where("name like ?", name+"%")
	}

	if tx.Find(&apps); tx.Error != nil {
		return nil, fmt.Errorf("repositories.sql.application.list: %w", tx.Error)
	}

	return apps, nil
}

func (sql *SqlApplication) Update(ctx context.Context, app *entities.Application) (*entities.Application, error) {
	app.UpdatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Model(app).Select("name", "updated_at").Updates(app); tx.Error != nil {
		return nil, fmt.Errorf("repositories.sql.application.create: %w", tx.Error)
	}

	return app, nil
}

func (sql *SqlApplication) Delete(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Non-repeatable_reads
	tx := sql.client.Begin(&xsql.TxOptions{Isolation: xsql.LevelReadCommitted})

	if txn := tx.Model(&app).Where("id = ?", id).First(&app); txn.Error != nil {
		return nil, fmt.Errorf("repositories.sql.application.delete.get: %w", txn.Error)
	}

	app.UpdatedAt = sql.timer.Now().UnixMilli()
	app.DeletedAt = sql.timer.Now().UnixMilli()

	if txn := tx.Model(app).Select("updated_at", "deleted_at").Updates(app); txn.Error != nil {
		return nil, fmt.Errorf("repositories.sql.application.delete.update: %w", txn.Error)
	}

	if txn := tx.Commit(); txn.Error != nil {
		return nil, fmt.Errorf("repositories.sql.application.delete: %w", tx.Error)
	}

	return &app, nil
}
