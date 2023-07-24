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

type SqlEndpointRule struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlEndpointRule) Create(ctx context.Context, entity *entities.EndpointRule) (*entities.EndpointRule, error) {
	entity.GenId()
	entity.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entity); tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.create: %w", tx.Error)
	}
	return entity, nil
}

func (sql *SqlEndpointRule) BulkCreate(ctx context.Context, entities []entities.EndpointRule) ([]string, error) {
	ids := []string{}
	for i, entity := range entities {
		entity.GenId()
		entity.SetAT(sql.timer.Now())

		ids = append(ids, entity.Id)
		entities[i] = entity
	}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(entities); tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.bulk_create: %w", tx.Error)
	}
	return ids, nil
}

func (sql *SqlEndpointRule) List(ctx context.Context, wsId, appId, epId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error) {
	req := structure.ListReqBuild(opts)

	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(epr).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."endpoint_id"`, ep.TableName(), ep.TableName(), epr.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), epId)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, epr.TableName()))

	res := &structure.ListRes[entities.EndpointRule]{Data: []entities.EndpointRule{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.list: %w", tx.Error)
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpointRule) Get(ctx context.Context, wsId, appId, epId, id string) (*entities.EndpointRule, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)

	ws := &entities.Workspace{}
	app := &entities.Application{}
	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}

	tx := transaction.WithContext(ctx).Model(epr).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."endpoint_id"`, ep.TableName(), ep.TableName(), epr.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`, app.TableName(), app.TableName(), ep.TableName())).
		Joins(fmt.Sprintf(`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`, ws.TableName(), ws.TableName(), app.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ws.TableName()), wsId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, app.TableName()), appId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), epId).
		Where(fmt.Sprintf(`"%s"."id" = ?`, epr.TableName()), id).
		First(epr)
	if err := database.ErrGet(tx); err != nil {
		return nil, fmt.Errorf("endpoint_rule.get: %w", err)
	}

	return epr, nil
}
