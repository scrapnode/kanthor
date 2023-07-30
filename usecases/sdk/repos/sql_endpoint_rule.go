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

func (sql *SqlEndpointRule) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error) {
	req := structure.ListReqBuild(opts)

	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(epr).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep)).
		Scopes(JoinEp(epr))
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, epr.TableName()))

	res := &structure.ListRes[entities.EndpointRule]{Data: []entities.EndpointRule{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpointRule) Get(ctx context.Context, wsId, id string) (*entities.EndpointRule, error) {
	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}
	epr.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(epr).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep)).
		Scopes(JoinEp(epr)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, epr.TableName()), id).
		First(epr)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return epr, nil
}

func (sql *SqlEndpointRule) Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	doc.GenId()
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Update(ctx context.Context, wsId string, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}
	doc.SetAT(sql.timer.Now())

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep)).
		Scopes(JoinEp(epr)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Delete(ctx context.Context, wsId, id string) error {
	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}
	epr.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(epr).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep)).
		Scopes(JoinEp(epr)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, epr.TableName()), id).
		Delete(epr)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
