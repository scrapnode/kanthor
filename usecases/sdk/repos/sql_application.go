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

func (sql *SqlApplication) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error) {
	req := structure.ListReqBuild(opts)

	ws := &entities.Workspace{}
	app := &entities.Application{}

	tx := sql.client.WithContext(ctx).Model(app).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId)
	tx = database.SqlToListQuery(tx, req, `"id"`)

	res := &structure.ListRes[entities.Application]{Data: []entities.Application{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlApplication) Get(ctx context.Context, wsId, id string) (*entities.Application, error) {
	ws := &entities.Workspace{}
	app := &entities.Application{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&app).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), id).
		First(app)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return app, nil
}

func (sql *SqlApplication) Create(ctx context.Context, doc *entities.Application) (*entities.Application, error) {
	doc.GenId()
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlApplication) Update(ctx context.Context, doc *entities.Application) (*entities.Application, error) {
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Updates(doc); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return doc, nil
}

func (sql *SqlApplication) Delete(ctx context.Context, wsId, id string) error {
	ws := &entities.Workspace{}
	app := &entities.Application{}
	app.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), id).
		Delete(app)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
