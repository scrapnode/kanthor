package repositories

import (
	"context"
	xsql "database/sql"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
	"strings"
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

func (sql *SqlEndpoint) List(ctx context.Context, appId, name string) ([]entities.Endpoint, error) {
	var endpoints []entities.Endpoint

	var tx = sql.client.Model(&entities.Endpoint{}).
		Scopes(NotDeleted(sql.timer, &entities.Endpoint{})).
		Where("app_id = ?", appId).
		Order("priority DESC, name ASC")

	if name != "" {
		tx = tx.Where("name like ?", name+"%")
	}

	if tx.Find(&endpoints); tx.Error != nil {
		return nil, fmt.Errorf("endpoint.list: %w", tx.Error)
	}

	return endpoints, nil
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

func (sql *SqlEndpoint) ListWithRules(ctx context.Context, appId string) ([]EndpointWithRules, error) {
	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}
	join := fmt.Sprintf("LEFT JOIN %s ON %s.endpoint_id = %s.id", epr.TableName(), epr.TableName(), ep.TableName())
	selects := strings.Join([]string{
		fmt.Sprintf("%s.id AS endpoint_id", ep.TableName()),
		fmt.Sprintf("%s.app_id AS endpoint_app_id", ep.TableName()),
		fmt.Sprintf("%s.name AS endpoint_name", ep.TableName()),
		fmt.Sprintf("%s.uri AS endpoint_uri", ep.TableName()),
		fmt.Sprintf("%s.method AS endpoint_method", ep.TableName()),
		fmt.Sprintf("%s.id AS rule_id", epr.TableName()),
		fmt.Sprintf("%s.priority AS rule_priority", epr.TableName()),
		fmt.Sprintf("%s.exclusionary AS rule_exclusionary", epr.TableName()),
		fmt.Sprintf("%s.condition_source AS rule_condition_source", epr.TableName()),
		fmt.Sprintf("%s.condition_expression AS rule_condition_expression", epr.TableName()),
	}, ",")
	rows, err := sql.client.
		Model(ep).
		Joins(join).
		Scopes(NotDeleted(sql.timer, ep)).
		Scopes(NotDeleted(sql.timer, epr)).
		Select(selects).
		Rows()
	if err != nil {
		return nil, err
	}

	maps := map[string]*EndpointWithRules{}
	for rows.Next() {
		var entity endpointWithRule
		if err := sql.client.ScanRows(rows, &entity); err != nil {
			return nil, err
		}

		if _, ok := maps[entity.EndpointId]; !ok {
			maps[entity.EndpointId] = &EndpointWithRules{
				Endpoint: entities.Endpoint{
					Entity: entities.Entity{Id: entity.EndpointId},
					AppId:  entity.EndpointAppId,
					Name:   entity.EndpointName,
					Uri:    entity.EndpointUri,
					Method: entity.EndpointMethod,
				},
				Rules: []entities.EndpointRule{},
			}
		}

		maps[entity.EndpointId].Rules = append(
			maps[entity.EndpointId].Rules,
			entities.EndpointRule{
				Entity:              entities.Entity{Id: entity.RuleId},
				Priority:            entity.RulePriority,
				Exclusionary:        entity.RuleExclusionary,
				ConditionSource:     entity.RuleConditionSource,
				ConditionExpression: entity.RuleConditionExpression,
			},
		)

	}

	var endpoints []EndpointWithRules
	for _, m := range maps {
		endpoints = append(endpoints, *m)
	}
	return endpoints, nil
}

type endpointWithRule struct {
	EndpointId              string `json:"endpoint_id"`
	EndpointAppId           string `json:"endpoint_app_id"`
	EndpointName            string `json:"endpoint_name"`
	EndpointUri             string `json:"endpoint_uri"`
	EndpointMethod          string `json:"endpoint_method"`
	RuleId                  string `json:"rule_id"`
	RulePriority            int    `json:"rule_priority"`
	RuleExclusionary        bool   `json:"rule_exclusionary"`
	RuleConditionSource     string `json:"rule_condition_source"`
	RuleConditionExpression string `json:"rule_condition_expression"`
}
