package repositories

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
	ws := &entities.Workspace{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&ws).
		Where(fmt.Sprintf(`"%s"."id" = ?`, entities.TableWs), id).
		First(ws)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return ws, nil
}
