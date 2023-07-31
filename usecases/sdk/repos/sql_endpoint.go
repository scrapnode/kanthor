package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	if tx := transaction.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Update(ctx context.Context, wsId string, doc *entities.Endpoint) (*entities.Endpoint, error) {
	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, wsId, id string) error {
	ep := &entities.Endpoint{}
	ep.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), id).
		Delete(ep)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpoint) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error) {
	ep := &entities.Endpoint{}
	tx := sql.client.WithContext(ctx).Model(ep).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, ep.TableName()))

	res := &structure.ListRes[entities.Endpoint]{Data: []entities.Endpoint{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId, id string) (*entities.Endpoint, error) {
	ep := &entities.Endpoint{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(ep).
		Scopes(UseWsId(wsId)).
		Scopes(JoinApp(ep)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, ep.TableName()), id).
		First(ep)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ep, nil
}
