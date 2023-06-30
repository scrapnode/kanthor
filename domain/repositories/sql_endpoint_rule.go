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

type SqlEndpointRule struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlEndpointRule) Create(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error) {
	epr.Id = utils.ID("epr")
	epr.CreatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Create(epr); tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.create: %w", tx.Error)
	}
	return epr, nil
}

func (sql *SqlEndpointRule) Get(ctx context.Context, id string) (*entities.EndpointRule, error) {
	var rule entities.EndpointRule
	if tx := sql.client.Model(&rule).Where("id = ?", id).First(&rule); tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.get: %w", tx.Error)
	}

	if rule.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("endpoint_rule.get.deleted: deleted_at:%d", rule.DeletedAt)
	}

	return &rule, nil
}

func (sql *SqlEndpointRule) List(ctx context.Context, epId string) ([]entities.EndpointRule, error) {
	var rules []entities.EndpointRule
	var tx = sql.client.Model(&entities.EndpointRule{}).
		Scopes(NotDeleted(sql.timer, &entities.EndpointRule{})).
		Where("endpoint_id = ?", epId).
		// example:
		// priority - exclusionary
		// 		  9 - TRUE
		// 		  9 - FALSE
		// 		  8 - FALSE
		// 		  7 - TRUE
		// 		  7 - FALSE
		Order("priority DESC, exclusionary DESC")

	if tx.Find(&rules); tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.list: %w", tx.Error)
	}

	return rules, nil
}

func (sql *SqlEndpointRule) Update(ctx context.Context, epr *entities.EndpointRule) (*entities.EndpointRule, error) {
	epr.UpdatedAt = sql.timer.Now().UnixMilli()

	tx := sql.client.Model(epr).
		Select("condition", "priority", "exclusionary", "updated_at").
		Updates(epr)
	if tx.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.create: %w", tx.Error)
	}

	return epr, nil
}

func (sql *SqlEndpointRule) Delete(ctx context.Context, id string) (*entities.EndpointRule, error) {
	var epr entities.EndpointRule
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Non-repreatable_reads
	tx := sql.client.Begin(&xsql.TxOptions{Isolation: xsql.LevelReadCommitted})

	if txn := tx.Model(&epr).Where("id = ?", id).First(&epr); txn.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.delete.get: %w", txn.Error)
	}

	epr.UpdatedAt = sql.timer.Now().UnixMilli()
	epr.DeletedAt = sql.timer.Now().UnixMilli()

	if txn := tx.Model(epr).Select("updated_at", "deleted_at").Updates(epr); txn.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.delete.update: %w", txn.Error)
	}

	if txn := tx.Commit(); txn.Error != nil {
		return nil, fmt.Errorf("endpoint_rule.delete: %w", tx.Error)
	}

	return &epr, nil
}
