package entities

type Application struct {
	Entity
	AuditTime
	// @TODO: add deactivated_at
	// DeactivatedAt int64 `json:"deactivated_at"`

	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
}

func (entity *Application) TableName() string {
	return "kanthor_application"
}

func (entity *Application) GenId() {
	if entity.Id == "" {
		entity.Id = AppId()
	}
}
