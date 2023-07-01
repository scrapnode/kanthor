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

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
}

func (sql *SqlApplication) Create(ctx context.Context, app *entities.Application) (*entities.Application, error) {
	app.Id = utils.ID("app")
	app.CreatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Create(app); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}
	return app, nil
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	if tx := sql.client.Model(&app).Where("id = ?", id).First(&app); tx.Error != nil {
		return nil, fmt.Errorf("application.get: %w", tx.Error)
	}

	if app.DeletedAt >= sql.timer.Now().UnixMilli() {
		return nil, fmt.Errorf("application.get.deleted: deleted_at:%d", app.DeletedAt)
	}

	return &app, nil
}

func (sql *SqlApplication) List(ctx context.Context, wsId, name string) ([]entities.Application, error) {
	var apps []entities.Application
	var tx = sql.client.Model(&entities.Application{}).
		Scopes(NotDeleted(sql.timer, &entities.Application{})).
		Where("workspace_id = ?", wsId)

	if name != "" {
		tx = tx.Where("name like ?", name+"%")
	}

	if tx.Find(&apps); tx.Error != nil {
		return nil, fmt.Errorf("application.list: %w", tx.Error)
	}

	return apps, nil
}

func (sql *SqlApplication) Update(ctx context.Context, app *entities.Application) (*entities.Application, error) {
	app.UpdatedAt = sql.timer.Now().UnixMilli()

	if tx := sql.client.Model(app).Select("name", "updated_at").Updates(app); tx.Error != nil {
		return nil, fmt.Errorf("application.create: %w", tx.Error)
	}

	return app, nil
}

func (sql *SqlApplication) Delete(ctx context.Context, id string) (*entities.Application, error) {
	var app entities.Application
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Non-repeatable_reads
	tx := sql.client.Begin(&xsql.TxOptions{Isolation: xsql.LevelReadCommitted})

	if txn := tx.Model(&app).Where("id = ?", id).First(&app); txn.Error != nil {
		return nil, fmt.Errorf("application.delete.get: %w", txn.Error)
	}

	app.UpdatedAt = sql.timer.Now().UnixMilli()
	app.DeletedAt = sql.timer.Now().UnixMilli()

	if txn := tx.Model(app).Select("updated_at", "deleted_at").Updates(app); txn.Error != nil {
		return nil, fmt.Errorf("application.delete.update: %w", txn.Error)
	}

	if txn := tx.Commit(); txn.Error != nil {
		return nil, fmt.Errorf("application.delete: %w", tx.Error)
	}

	return &app, nil
}

func (sql *SqlApplication) GetWithWorkspace(ctx context.Context, id string) (*ApplicationWithWorkspace, error) {
	app := &entities.Application{}
	ws := &entities.Workspace{}
	wst := &entities.WorkspaceTier{}

	appws := fmt.Sprintf("JOIN %s ON %s.id = %s.workspace_id", ws.TableName(), ws.TableName(), app.TableName())
	wswst := fmt.Sprintf("JOIN %s ON %s.workspace_id = %s.id", wst.TableName(), wst.TableName(), ws.TableName())
	selects := strings.Join([]string{
		fmt.Sprintf("%s.id AS app_id", app.TableName()),
		fmt.Sprintf("%s.name AS app_name", app.TableName()),
		fmt.Sprintf("%s.id AS ws_id", ws.TableName()),
		fmt.Sprintf("%s.owner_id AS ws_owner_id", ws.TableName()),
		fmt.Sprintf("%s.name AS ws_name", ws.TableName()),
		fmt.Sprintf("%s.name AS ws_tier_name", wst.TableName()),
	}, ", ")

	var entity applicationWithWorkspace
	tx := sql.client.
		Model(app).
		Joins(appws).
		Joins(wswst).
		Scopes(NotDeleted(sql.timer, app)).
		Scopes(NotDeleted(sql.timer, ws)).
		Select(selects).
		First(&entity)
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := ApplicationWithWorkspace{
		Application: entities.Application{
			Entity:      entities.Entity{Id: entity.AppId},
			WorkspaceId: entity.WsId,
			Name:        entity.AppName,
		},
		Workspace: entities.Workspace{
			Entity:  entities.Entity{Id: entity.WsId},
			OwnerId: entity.WsOwnerId,
			Name:    entity.WsName,
			Tier:    &entities.WorkspaceTier{WorkspaceId: entity.WsId, Name: entity.WsTierName},
		},
	}

	return &result, nil
}

type applicationWithWorkspace struct {
	AppId      string `json:"app_id"`
	AppName    string `json:"app_name"`
	WsId       string `json:"ws_id"`
	WsOwnerId  string `json:"ws_owner_id"`
	WsName     string `json:"ws_name"`
	WsTierName string `json:"ws_tier_name"`
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

	result := ApplicationWithEndpointsAndRules{
		Application: *app,
		Endpoints:   []EndpointWithRules{},
	}
	for _, m := range maps {
		result.Endpoints = append(result.Endpoints, *m)
	}
	return &result, nil
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
