package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) CreateBulk(ctx context.Context, docs []entities.Endpoint) ([]string, error) {
	ids := []string{}
	if len(docs) == 0 {
		return ids, nil
	}

	for _, doc := range docs {
		ids = append(ids, doc.Id)
	}

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(docs); tx.Error != nil {
		return nil, tx.Error
	}
	return ids, nil
}

func (sql *SqlEndpoint) Count(ctx context.Context, wsId string, query *entities.PagingQuery) (int64, error) {
	doc := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseApp(doc.TableName()),
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

func (sql *SqlEndpoint) Get(ctx context.Context, wsId, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}
	doc.Id = id

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&doc).
		Scopes(UseApp(doc.TableName()), UseWsId(wsId, entities.TableApp)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}
