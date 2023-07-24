package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
	"strings"
)

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	if tx := sql.client.Model(&app).Where("id = ?", id).First(&app); tx.Error != nil {
		return nil, fmt.Errorf("application.get: %w", tx.Error)
	}

	return &app, nil
}

func (sql *SqlApplication) ListEndpointsWithRules(ctx context.Context, id string) (*ApplicationWithEndpointsAndRules, error) {
	app, err := sql.Get(ctx, id)
	if err != nil {
		return nil, err
	}

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
	}, ", ")
	rows, err := sql.client.
		Model(ep).
		Joins(join).
		Where(fmt.Sprintf("%s.app_id = ?", ep.TableName()), app.Id).
		Select(selects).
		Rows()
	if err != nil {
		return nil, err
	}

	maps := map[string]*EndpointWithRules{}
	for rows.Next() {
		var entity epwr
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

	result := ApplicationWithEndpointsAndRules{
		Application: *app,
		Endpoints:   []EndpointWithRules{},
	}
	for _, m := range maps {
		result.Endpoints = append(result.Endpoints, *m)
	}
	return &result, nil
}

type epwr struct {
	EndpointId              string `json:"endpoint_id"`
	EndpointAppId           string `json:"endpoint_app_id"`
	EndpointName            string `json:"endpoint_name"`
	EndpointUri             string `json:"endpoint_uri"`
	EndpointMethod          string `json:"endpoint_method"`
	RuleId                  string `json:"rule_id"`
	RulePriority            int32  `json:"rule_priority"`
	RuleExclusionary        bool   `json:"rule_exclusionary"`
	RuleConditionSource     string `json:"rule_condition_source"`
	RuleConditionExpression string `json:"rule_condition_expression"`
}
