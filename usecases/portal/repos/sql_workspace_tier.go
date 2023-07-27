package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlWorkspaceTier struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspaceTier) Create(ctx context.Context, doc *entities.WorkspaceTier) (*entities.WorkspaceTier, error) {
	doc.GenId()
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspaceTier) BulkCreate(ctx context.Context, docs []entities.WorkspaceTier) ([]string, error) {
	ids := []string{}
	if len(docs) == 0 {
		return ids, nil
	}

	now := sql.timer.Now()
	for i, doc := range docs {
		doc.GenId()
		doc.SetAT(now)

		ids = append(ids, doc.Id)
		docs[i] = doc
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(docs); tx.Error != nil {
		return nil, tx.Error
	}

	return ids, nil
}
