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

func (sql *SqlEndpointRule) Update(ctx context.Context, wsId string, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(&entities.Endpoint{})).
		Scopes(JoinEp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Delete(ctx context.Context, wsId, id string) error {
	doc := &entities.EndpointRule{}
	doc.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(&entities.Endpoint{})).
		Scopes(JoinEp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpointRule) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.EndpointRule], error) {
	doc := &entities.EndpointRule{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(&entities.Endpoint{})).
		Scopes(JoinEp(doc))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, doc.TableName()))

	res := &structure.ListRes[entities.EndpointRule]{Data: []entities.EndpointRule{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpointRule) Get(ctx context.Context, wsId, id string) (*entities.EndpointRule, error) {
	doc := &entities.EndpointRule{}
	doc.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(&entities.Endpoint{})).
		Scopes(JoinEp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
