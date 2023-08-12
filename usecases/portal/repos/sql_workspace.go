package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
}

func (sql *SqlWorkspace) BulkCreate(ctx context.Context, docs []entities.Workspace) ([]string, error) {
	ids := []string{}
	if len(docs) == 0 {
		return ids, nil
	}

	for i, doc := range docs {
		ids = append(ids, doc.Id)
		docs[i] = doc
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(docs); tx.Error != nil {
		return nil, tx.Error
	}

	return ids, nil
}

func (sql *SqlWorkspace) Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspace) Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlWorkspace) List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error) {
	req := structure.ListReqBuild(opts)
	doc := &entities.Workspace{}

	tx := sql.client.WithContext(ctx).Model(doc)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, doc.TableName()))

	res := &structure.ListRes[entities.Workspace]{Data: []entities.Workspace{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("workspace.list: %w", tx.Error)
	}

	return res, nil
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}

	tx := sql.client.WithContext(ctx).Model(&doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspace) GetOwned(ctx context.Context, owner string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}
	tx := sql.client.WithContext(ctx).Model(&doc).
		Where("owner_id = ?", owner).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
