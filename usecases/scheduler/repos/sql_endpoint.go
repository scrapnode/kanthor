package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) List(ctx context.Context, appId string) ([]entities.Endpoint, error) {
	docs := []entities.Endpoint{}
	tx := sql.client.WithContext(ctx).Where("app_id = ?", appId).Find(&docs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return docs, nil
}
