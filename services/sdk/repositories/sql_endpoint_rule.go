package repositories

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlEndpointRule struct {
	client *gorm.DB
}

func (sql *SqlEndpointRule) Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Update(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Delete(ctx context.Context, doc *entities.EndpointRule) error {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpointRule) List(ctx context.Context, wsId, appId, epId string, query *entities.Query) ([]entities.EndpointRule, error) {
	doc := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEpId(epId, doc.TableName()),
			UseAppId(appId, entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		).
		Order(clause.OrderByColumn{Column: clause.Column{Name: fmt.Sprintf(`"%s"."created_at"`, doc.TableName())}, Desc: true})

	if len(query.Ids) > 0 {
		tx = tx.Where(fmt.Sprintf(`"%s"."id" IN ?`, doc.TableName()), query.Ids)
	} else {
		props := []string{
			fmt.Sprintf(`"%s"."name"`, doc.TableName()),
			fmt.Sprintf(`"%s"."condition_source"`, doc.TableName()),
		}
		tx = database.ApplyListQuery(tx, props, query.Search, query.Limit, query.Page)
	}

	var docs []entities.EndpointRule
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlEndpointRule) Count(ctx context.Context, wsId, appId, epId string, query *entities.Query) (int64, error) {
	doc := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEpId(epId, doc.TableName()),
			UseAppId(appId, entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		)

	if len(query.Ids) > 0 {
		return int64(len(query.Ids)), nil
	}

	props := []string{
		fmt.Sprintf(`"%s"."name"`, doc.TableName()),
		fmt.Sprintf(`"%s"."condition_source"`, doc.TableName()),
	}
	tx = database.ApplyCountQuery(tx, props, query.Search)
	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlEndpointRule) Get(ctx context.Context, wsId string, id string) (*entities.EndpointRule, error) {
	doc := &entities.EndpointRule{}
	doc.Id = id

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(
			UseEp(doc.TableName()),
			UseApp(entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}
