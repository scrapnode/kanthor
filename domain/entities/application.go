package entities

import "encoding/json"

type Application struct {
	Entity
	AuditTime
	// @TODO: add deactivated_at
	// DeactivatedAt int64 `json:"deactivated_at"`

	WsId string `json:"ws_id"`
	Name string `json:"name"`
}

func (entity *Application) TableName() string {
	return TableApp
}

func (entity *Application) GenId() {
	if entity.Id == "" {
		entity.Id = AppId()
	}
}

func (entity *Application) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Application) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *Application) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}

type ApplicationWithRelationship struct {
	*Application
	Workspace *Workspace
}
