package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Create(ctx context.Context, doc *entities.Application) (*entities.Application, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlApplication) Update(ctx context.Context, wsId string, doc *entities.Application) (*entities.Application, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(wsId)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)

	if tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return doc, nil
}

func (sql *SqlApplication) Delete(ctx context.Context, wsId, id string) error {
	app := &entities.Application{}
	app.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(wsId)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), id).
		Delete(app)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sql *SqlApplication) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error) {
	app := &entities.Application{}
	tx := sql.client.WithContext(ctx).Model(app).
		Scopes(UseWsId(wsId))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, `"id"`)

	res := &structure.ListRes[entities.Application]{Data: []entities.Application{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlApplication) Get(ctx context.Context, wsId, id string) (*entities.Application, error) {
	app := &entities.Application{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&app).
		Scopes(UseWsId(wsId)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), id).
		First(app)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return app, nil
}
