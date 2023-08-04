package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type SqlWorkspaceTier struct {
	client *gorm.DB
}

func (sql *SqlWorkspaceTier) Get(ctx context.Context, wsId string) (*entities.WorkspaceTier, error) {
	doc := &entities.WorkspaceTier{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc)).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
