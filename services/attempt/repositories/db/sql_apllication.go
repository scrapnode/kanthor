package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.ApplicationWithRelationship, error) {
	doc := &entities.Application{}

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	returning := &entities.ApplicationWithRelationship{Application: doc}
	wstx := transaction.WithContext(ctx).Model(&entities.Workspace{}).
		Where(fmt.Sprintf(`"%s"."id" = ?`, entities.TableWs), doc.WsId).
		First(&returning.Workspace)
	if wstx.Error != nil {
		return nil, database.SqlError(wstx.Error)
	}

	return returning, nil
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
