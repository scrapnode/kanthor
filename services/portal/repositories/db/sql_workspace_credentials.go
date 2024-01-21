package db

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	updates := []string{
		"name",
		"expired_at",
		"updated_at",
	}
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		// When update with struct, GORM will only update non-zero fields,
		// you might want to use map to update attributes or use Select to specify fields to update
		Select(updates).
		Updates(doc)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlWorkspaceCredentials) List(ctx context.Context, wsId string, query *entities.PagingQuery) ([]entities.WorkspaceCredentials, error) {
	doc := &entities.WorkspaceCredentials{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Order(clause.OrderByColumn{Column: clause.Column{Name: fmt.Sprintf(`"%s"."created_at"`, doc.TableName())}, Desc: true})

	if len(query.Ids) > 0 {
		tx = tx.Where(fmt.Sprintf(`"%s"."id" IN ?`, doc.TableName()), query.Ids)
	} else {
		props := []string{fmt.Sprintf(`"%s"."name"`, doc.TableName())}
		tx = database.SqlApplyListQuery(tx, query, props)
	}

	var docs []entities.WorkspaceCredentials
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlWorkspaceCredentials) Count(ctx context.Context, wsId string, query *entities.PagingQuery) (int64, error) {
	doc := &entities.WorkspaceCredentials{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, entities.TableWsc))

	if len(query.Ids) > 0 {
		return int64(len(query.Ids)), nil
	}

	props := []string{fmt.Sprintf(`"%s"."name"`, doc.TableName())}
	tx = database.SqlApplyCountQuery(tx, query, props)
	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlWorkspaceCredentials) Get(ctx context.Context, wsId, id string) (*entities.WorkspaceCredentials, error) {
	doc := &entities.WorkspaceCredentials{}
	doc.Id = id

	transaction := database.SqlTxnFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
