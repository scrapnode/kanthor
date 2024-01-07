package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}
	transaction := database.SqlTxnFromContext(ctx, sql.client)

	tx := transaction.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}
