package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpointRule struct {
	client *gorm.DB
}

func (sql *SqlEndpointRule) CreateBulk(ctx context.Context, docs []entities.EndpointRule) ([]string, error) {
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

func (sql *SqlEndpointRule) Count(ctx context.Context, wsId string, query *entities.PagingQuery) (int64, error) {
	doc := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEp(doc.TableName()),
			UseApp(entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		)

	if len(query.Ids) > 0 {
		return int64(len(query.Ids)), nil
	}

	props := []string{
		fmt.Sprintf(`"%s"."name"`, doc.TableName()),
		fmt.Sprintf(`"%s"."condition_source"`, doc.TableName()),
	}
	tx = database.SqlApplyListQuery(tx, query, props)
	var count int64
	return count, tx.Count(&count).Error
}
