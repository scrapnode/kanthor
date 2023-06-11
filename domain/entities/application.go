package entities

type Application struct {
	Entity
	AuditTime
	SoftDelete

	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
}

func (entity *Application) TableName() string {
	return "application"
}
