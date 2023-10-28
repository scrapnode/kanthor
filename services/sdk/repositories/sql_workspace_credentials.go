package repositories

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlWorkspaceCredentials struct {
	client *gorm.DB
}

func (sql *SqlWorkspaceCredentials) Get(ctx context.Context, id string) (*entities.WorkspaceCredentials, error) {
	wsc := &entities.WorkspaceCredentials{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(wsc).
		Where(fmt.Sprintf(`"%s".id = ?`, entities.TableWsc), id).
		First(wsc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return wsc, nil
}
