package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) BulkCreate(ctx context.Context, docs []entities.Application) ([]string, error) {
	ids := []string{}
	if len(docs) == 0 {
		return ids, nil
	}

	for i, doc := range docs {
		ids = append(ids, doc.Id)
		docs[i] = doc
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(docs); tx.Error != nil {
		return nil, tx.Error
	}
	return ids, nil
}
