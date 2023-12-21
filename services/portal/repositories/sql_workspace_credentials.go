package repositories

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlWorkspaceCredentials struct {
	client *gorm.DB
}

func (sql *SqlWorkspaceCredentials) Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspaceCredentials) Update(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error) {
	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlWorkspaceCredentials) List(ctx context.Context, wsId string, q string, limit, page int) ([]entities.WorkspaceCredentials, error) {
	doc := &entities.WorkspaceCredentials{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Order(fmt.Sprintf("%s.id DESC", doc.TableName()))

	tx = database.ApplyListQuery(tx, q, []string{fmt.Sprintf("%s.name", doc.TableName())}, limit, page)

	var docs []entities.WorkspaceCredentials
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlWorkspaceCredentials) Count(ctx context.Context, wsId, q string) (int64, error) {
	doc := &entities.WorkspaceCredentials{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, entities.TableWsc))

	tx = database.ApplyCountQuery(tx, q, []string{fmt.Sprintf("%s.name", doc.TableName())})

	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlWorkspaceCredentials) Get(ctx context.Context, wsId, id string) (*entities.WorkspaceCredentials, error) {
	doc := &entities.WorkspaceCredentials{}
	doc.Id = id

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(&doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, database.SqlError(tx.Error)
	}

	return doc, nil
}
