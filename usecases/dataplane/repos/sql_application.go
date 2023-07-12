package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
	"strings"
)

type SqlApplication struct {
	client *gorm.DB
	timer  timer.Timer
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

	var entity apws
	tx := sql.client.
		Model(app).
		Joins(appws).
		Joins(wswst).
		Scopes(database.NotDeleted(sql.timer, app)).
		Scopes(database.NotDeleted(sql.timer, ws)).
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

type apws struct {
	AppId      string `json:"app_id"`
	AppName    string `json:"app_name"`
	WsId       string `json:"ws_id"`
	WsOwnerId  string `json:"ws_owner_id"`
	WsName     string `json:"ws_name"`
	WsTierName string `json:"ws_tier_name"`
}
