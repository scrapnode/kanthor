package entities

import "github.com/scrapnode/kanthor/pkg/utils"

type Application struct {
	Entity
	AuditTime
	SoftDelete

	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
}

func (entity *Application) TableName() string {
	return "kanthor_application"
}

func (entity *Application) GenId() {
	entity.Id = utils.ID("app")
}
