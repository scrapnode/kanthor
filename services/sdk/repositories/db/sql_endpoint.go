package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)

	updates := []string{
		"name",
		"method",
		"uri",
		"updated_at",
	}
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		// When update with struct, GORM will only update non-zero fields,
		// you might want to use map to update attributes or use Select to specify fields to update
		Select(updates).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, doc *entities.Endpoint) error {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpoint) List(ctx context.Context, wsId, appId string, query *entities.PagingQuery) ([]entities.Endpoint, error) {
	doc := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseAppId(appId, doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		).
		Order(clause.OrderByColumn{Column: clause.Column{Name: fmt.Sprintf(`"%s"."created_at"`, doc.TableName())}, Desc: true})

	if len(query.Ids) > 0 {
		tx = tx.Where(fmt.Sprintf(`"%s"."id" IN ?`, doc.TableName()), query.Ids)
	} else {
		props := []string{
			fmt.Sprintf(`"%s"."name"`, doc.TableName()),
			fmt.Sprintf(`"%s"."uri"`, doc.TableName()),
		}
		tx = database.SqlApplyListQuery(tx, query, props)
	}

	var docs []entities.Endpoint
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlEndpoint) Count(ctx context.Context, wsId, appId string, query *entities.PagingQuery) (int64, error) {
	doc := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseAppId(appId, doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		)

	if len(query.Ids) > 0 {
		return int64(len(query.Ids)), nil
	}

	props := []string{
		fmt.Sprintf(`"%s"."name"`, doc.TableName()),
		fmt.Sprintf(`"%s"."uri"`, doc.TableName()),
	}
	tx = database.SqlApplyListQuery(tx, query, props)
	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId string, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}
	doc.Id = id

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(
			UseApp(doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
