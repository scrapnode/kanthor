package repositories

import (
	"context"
	xsql "database/sql"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlEndpoint) Create(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error) {
	ep.Id = utils.ID("ep")
	ep.CreatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Create(ep); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.create: %w", tx.Error)
	}
	return ep, nil
}

func (sql *SqlEndpoint) Get(ctx context.Context, id string) (*entities.Endpoint, error) {
	var ep entities.Endpoint
	if tx := sql.client.Model(&ep).Where("id = ?", id).First(&ep); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.get: %w", tx.Error)
	}

	if ep.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("endpoint.get.deleted: deleted_at:%d", ep.DeletedAt)
	}

	return &ep, nil
}

func (sql *SqlEndpoint) List(ctx context.Context, appId string, opts ...ListOps) (*ListRes[entities.Endpoint], error) {
	req := ListReqBuild(opts)
	res := ListRes[entities.Endpoint]{}

	var tx = sql.client.Model(&entities.Endpoint{}).
		Scopes(NotDeleted(sql.timer, &entities.Endpoint{})).
		Where("app_id = ?", appId).
		Order("priority DESC, name ASC")

	tx = TxListCursor(tx, req)
	if req.Search != "" {
		tx = tx.Where("name like ?", req.Search+"%")
	}

	if tx.Find(&res.Data); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.list: %w", tx.Error)
	}

	return &res, nil
}

func (sql *SqlEndpoint) Update(ctx context.Context, ep *entities.Endpoint) (*entities.Endpoint, error) {
	ep.UpdatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Model(ep).Select("name", "uri", "updated_at").Updates(ep); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.create: %w", tx.Error)
	}

	return ep, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, id string) (*entities.Endpoint, error) {
	var ep entities.Endpoint
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Non-repeatable_reads
	tx := sql.client.Begin(&xsql.TxOptions{Isolation: xsql.LevelReadCommitted})

	if txn := tx.Model(&ep).Where("id = ?", id).First(&ep); txn.Error != nil {
		return nil, fmt.Errorf("endpoint.delete.get: %w", txn.Error)
	}

	ep.UpdatedAt = sql.timer.Now().UnixMilli()
	ep.DeletedAt = sql.timer.Now().UnixMilli()

	if txn := tx.Model(ep).Select("updated_at", "deleted_at").Updates(ep); txn.Error != nil {
		return nil, fmt.Errorf("endpoint.delete.update: %w", txn.Error)
	}

	if txn := tx.Commit(); txn.Error != nil {
		return nil, fmt.Errorf("endpoint.delete: %w", tx.Error)
	}

	return &ep, nil
}
