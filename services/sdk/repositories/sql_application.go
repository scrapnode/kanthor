package repositories

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
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

func (sql *SqlApplication) List(ctx context.Context, wsId string, q string, limit, page int) ([]entities.Application, error) {
	doc := &entities.Application{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName()))

	tx = database.ApplyListQuery(tx, q, []string{fmt.Sprintf("%s.name", doc.TableName())}, limit, page)

	var docs []entities.Application
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlApplication) Count(ctx context.Context, wsId string, q string) (int64, error) {
	doc := &entities.Application{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, entities.TableWsc))

	tx = database.ApplyCountQuery(tx, q, []string{fmt.Sprintf("%s.name", doc.TableName())})

	var count int64
	return count, tx.Count(&count).Error
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
