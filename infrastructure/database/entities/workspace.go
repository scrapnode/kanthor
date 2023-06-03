package entities

import "github.com/scrapnode/kanthor/infrastructure/utils"

type Workspace struct {
	Entity
	AuditTime
	SoftDelete

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
}

func (entity *Workspace) GenId() {
	entity.Id = utils.ID("ws")
}
