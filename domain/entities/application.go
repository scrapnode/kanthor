package entities

import "github.com/scrapnode/kanthor/pkg/utils"

type Application struct {
	Entity
	AuditTime

	WorkspaceId string `json:"workspace_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
}

func (entity *Application) TableName() string {
	return "kanthor_application"
}

func (entity *Application) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("app")
	}
}
