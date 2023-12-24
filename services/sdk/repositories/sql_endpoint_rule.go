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

func (sql *SqlEndpointRule) List(ctx context.Context, wsId, appId, epId string, q string, limit, page int) ([]entities.EndpointRule, error) {
	doc := &entities.EndpointRule{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEpId(epId, doc.TableName()),
			UseAppId(appId, entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		).
		Order(clause.OrderByColumn{Column: clause.Column{Name: fmt.Sprintf("%s.created_at", doc.TableName())}, Desc: true})

	qcols := []string{
		fmt.Sprintf("%s.name", doc.TableName()),
		fmt.Sprintf("%s.condition_source", doc.TableName()),
	}
	tx = database.ApplyListQuery(tx, q, qcols, limit, page)

	var docs []entities.EndpointRule
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlEndpointRule) Count(ctx context.Context, wsId, appId, epId string, q string) (int64, error) {
	doc := &entities.EndpointRule{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEpId(epId, doc.TableName()),
			UseAppId(appId, entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		)

	qcols := []string{
		fmt.Sprintf("%s.name", doc.TableName()),
		fmt.Sprintf("%s.condition_source", doc.TableName()),
	}
	tx = database.ApplyCountQuery(tx, q, qcols)

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
