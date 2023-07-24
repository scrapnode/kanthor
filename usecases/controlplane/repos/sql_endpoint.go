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

type SqlEndpoint struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlEndpoint) Create(ctx context.Context, entity *entities.Endpoint) (*entities.Endpoint, error) {
	entity.GenId()
	entity.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entity); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.create: %w", tx.Error)
	}
	return entity, nil
}

func (sql *SqlEndpoint) BulkCreate(ctx context.Context, entities []entities.Endpoint) ([]string, error) {
	ids := []string{}
	for i, entity := range entities {
		entity.GenId()
		entity.SetAT(sql.timer.Now())

		ids = append(ids, entity.Id)
		entities[i] = entity
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entities); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.bulk_create: %w", tx.Error)
	}
	return ids, nil
}

func (sql *SqlEndpoint) List(ctx context.Context, wsId, appId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error) {
	req := structure.ListReqBuild(opts)

	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(ep).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, ep.TableName()))

	res := &structure.ListRes[entities.Endpoint]{Data: []entities.Endpoint{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.list: %w", tx.Error)
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId, appId, id string) (*entities.Endpoint, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}

	tx := transaction.WithContext(ctx).Model(ep).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), id).
		First(ep)
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("endpoint.get: %w", err)
	}

	return ep, nil
}
