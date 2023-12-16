package repositories

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/structure"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Create(ctx context.Context, doc *entities.Application) (*entities.Application, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlApplication) Update(ctx context.Context, doc *entities.Application) (*entities.Application, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlApplication) Delete(ctx context.Context, doc *entities.Application) error {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sql *SqlApplication) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Application], error) {
	doc := &entities.Application{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName()))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, doc.TableName()))

	res := &structure.ListRes[entities.Application]{Data: []entities.Application{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlApplication) Get(ctx context.Context, wsId, id string) (*entities.Application, error) {
	doc := &entities.Application{}
	doc.Id = id

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}
