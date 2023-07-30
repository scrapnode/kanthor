package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlWorkspaceCredentials struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspaceCredentials) Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error) {
	doc.GenId()
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspaceCredentials) BulkCreate(ctx context.Context, docs []entities.WorkspaceCredentials) ([]string, error) {
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
