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

func (sql *SqlWorkspace) Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	if err := doc.Validate(); err != nil {
		return nil, err
	}

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspace) Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	if err := doc.Validate(); err != nil {
		return nil, err
	}

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	return doc, tx.Error
}

func (sql *SqlWorkspace) ListByIds(ctx context.Context, ids []string) ([]entities.Workspace, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var docs []entities.Workspace

	tx := sql.client.WithContext(ctx).
		Model(&entities.Workspace{}).
		Where("id IN ?", ids).
		Find(&docs)

	return docs, tx.Error
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}
	transaction := database.SqlTxnFromContext(ctx, sql.client)

	tx := transaction.WithContext(ctx).Model(&doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}

func (sql *SqlWorkspace) GetOwned(ctx context.Context, owner, id string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&doc).
		Where("owner_id = ? AND id = ?", owner, id).
		Order("id asc").
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}

func (sql *SqlWorkspace) ListOwned(ctx context.Context, owner string) ([]entities.Workspace, error) {
	var docs []entities.Workspace

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&entities.Workspace{}).
		Where("owner_id = ?", owner).
		Order("id DESC").
		Find(&docs)

	return docs, database.SqlError(tx.Error)
}
