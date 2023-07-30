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

type SqlEndpoint struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlEndpoint) List(ctx context.Context, wsId, appId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error) {
	req := structure.ListReqBuild(opts)

	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(ep).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, ep.TableName()))

	res := &structure.ListRes[entities.Endpoint]{Data: []entities.Endpoint{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId, appId, id string) (*entities.Endpoint, error) {
	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(ep).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), id).
		First(ep)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ep, nil
}

func (sql *SqlEndpoint) Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	doc.GenId()
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Updates(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, wsId, appId, id string) error {
	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}
	ep.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), id).
		Delete(ep)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
