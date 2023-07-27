package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlWorkspace) Create(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error) {
	doc.GenId()
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspace) BulkCreate(ctx context.Context, docs []entities.Workspace) ([]string, error) {
	ids := []string{}
	if len(docs) == 0 {
		return ids, nil
	}

	now := sql.timer.Now()
	for i, doc := range docs {
		doc.GenId()
		doc.SetAT(now)

		ids = append(ids, doc.Id)
		docs[i] = doc
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(docs); tx.Error != nil {
		return nil, tx.Error
	}

	return ids, nil
}

func (sql *SqlWorkspace) List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error) {
	req := structure.ListReqBuild(opts)
	ws := &entities.Workspace{}

	tx := sql.client.WithContext(ctx).Model(ws).
		Preload("Tier")
	tx = database.SqlToListQuery(tx, req, `"id"`)

	res := &structure.ListRes[entities.Workspace]{Data: []entities.Workspace{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("workspace.list: %w", tx.Error)
	}

	return res, nil
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	ws := &entities.Workspace{}

	tx := sql.client.WithContext(ctx).Model(&ws).
		Preload("Tier").
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), id).
		First(ws)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ws, nil
}

func (sql *SqlWorkspace) GetOwned(ctx context.Context, owner string) (*entities.Workspace, error) {
	ws := &entities.Workspace{}
	tx := sql.client.WithContext(ctx).Model(&ws).
		Preload("Tier").
		Where("owner_id = ?", owner).
		First(ws)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ws, nil
}
