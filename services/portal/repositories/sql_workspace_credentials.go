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

func (sql *SqlWorkspaceCredentials) List(ctx context.Context, wsId string, limit, page int, q string) ([]entities.WorkspaceCredentials, error) {
	tx := sql.client.WithContext(ctx).
		Model(&entities.WorkspaceCredentials{}).
		Scopes(UseWsId(wsId, entities.TableWsc)).
		Order(fmt.Sprintf("%s.id DESC", entities.TableWsc))

	tx = database.ApplyListQuery(tx, limit, page, q, []string{fmt.Sprintf("%s.name", entities.TableWsc)})

	var docs []entities.WorkspaceCredentials
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlWorkspaceCredentials) Count(ctx context.Context, wsId, q string) (int64, error) {
	tx := sql.client.WithContext(ctx).Model(&entities.WorkspaceCredentials{}).
		Scopes(UseWsId(wsId, entities.TableWsc))

	if q != "" {
		tx = tx.Where(fmt.Sprintf(`"%s.name" LIKE ?`, entities.TableWsc), "%"+q+"%")
	}

	var count int64
	tx = tx.Count(&count)

	return count, tx.Error
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
