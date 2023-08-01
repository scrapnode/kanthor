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
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, wsId, id string) error {
	doc := &entities.Endpoint{}
	doc.Id = id

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpoint) List(ctx context.Context, wsId string, opts ...structure.ListOps) (*structure.ListRes[entities.Endpoint], error) {
	doc := &entities.Endpoint{}
	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(doc))

	req := structure.ListReqBuild(opts)
	tx = database.SqlToListQuery(tx, req, fmt.Sprintf(`"%s"."id"`, doc.TableName()))

	res := &structure.ListRes[entities.Endpoint]{Data: []entities.Endpoint{}}
	if tx = tx.Find(&res.Data); tx.Error != nil {
		return nil, tx.Error
	}

	return structure.ListResBuild(res, req), nil
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}

	transaction := database.SqlClientFromContext(ctx, sql.client)
	tx := transaction.WithContext(ctx).Model(doc).
		Scopes(UseWsId(&entities.Application{}, wsId)).
		Scopes(JoinApp(doc)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
