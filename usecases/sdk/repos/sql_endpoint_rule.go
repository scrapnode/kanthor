package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type SqlEndpointRule struct {
	client *gorm.DB
}

func (sql *SqlEndpointRule) Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Update(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Delete(ctx context.Context, doc *entities.EndpointRule) error {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpointRule) List(ctx context.Context, wsId, appId, epId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error) {
	app := &entities.Application{}
	ep := &entities.Endpoint{}
	doc := &entities.EndpointRule{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseEpId(epId, doc)).
		Scopes(UseAppId(appId, ep)).
		Scopes(UseWsId(wsId, app))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, doc.TableName()))

	res := &structure.ListRes[entities.EndpointRule]{Data: []entities.EndpointRule{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpointRule) Get(ctx context.Context, wsId, appId, epId, id string) (*entities.EndpointRule, error) {
	app := &entities.Application{}
	ep := &entities.Endpoint{}
	doc := &entities.EndpointRule{}
	doc.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseEpId(epId, doc)).
		Scopes(UseAppId(appId, ep)).
		Scopes(UseWsId(wsId, app)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
