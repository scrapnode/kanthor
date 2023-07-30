package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	ws := &entities.Workspace{}

	tx := sql.client.WithContext(ctx).Model(&ws).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), id).
		First(ws)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ws, nil
}
