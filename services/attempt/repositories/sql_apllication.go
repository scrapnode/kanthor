package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Scan(ctx context.Context, limit int, cursor string) ([]entities.Application, error) {
	docs := []entities.Application{}
	tx := sql.client.WithContext(ctx).Limit(limit).Order("id ASC")
	if cursor != "" {
		tx = tx.Where("id > ?", cursor)
	}

	tx = tx.Find(&docs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return docs, nil
}

func (sql *SqlApplication) GetTiers(ctx context.Context, apps []entities.Application) (map[string]string, error) {
	var uniq = map[string]string{}
	for _, app := range apps {
		uniq[app.WsId] = app.WsId
	}

	docs := []entities.Workspace{}
	tx := sql.client.WithContext(ctx).Order("id ASC").Find(&docs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	tiers := map[string]string{}
	for _, doc := range docs {
		tiers[doc.Id] = doc.Tier
	}

	returning := map[string]string{}
	for _, app := range apps {
		returning[app.Id] = tiers[app.WsId]
	}

	return returning, nil
}
