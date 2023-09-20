package repos

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, doc *entities.Endpoint) error {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpoint) List(ctx context.Context, app *entities.Application, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error) {
	doc := &entities.Endpoint{}
	tx := sql.client.WithContext(ctx).Model(doc).Scopes(UseAppId(app.Id, doc))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, doc.TableName()))

	res := &structure.ListRes[entities.Endpoint]{Data: []entities.Endpoint{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpoint) Get(ctx context.Context, app *entities.Application, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}
	doc.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseAppId(app.Id, doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}

func (sql *SqlEndpoint) GetOfWorkspace(ctx context.Context, ws *entities.Workspace, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}
	doc.Id = id

	app := &entities.Application{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Joins(fmt.Sprintf(`JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), doc.TableName())).
		Joins(fmt.Sprintf(`JOIN "%s" ON "%s"."id" = "%s"."workspace_id" AND "%s" = ? `, ws.TableName(), ws.TableName(), app.TableName(), app.TableName()), ws.Id).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)

	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}
	return doc, nil
}
