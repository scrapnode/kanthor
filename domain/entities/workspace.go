package entities

import "github.com/scrapnode/kanthor/pkg/utils"

type Workspace struct {
	Entity
	AuditTime
	SoftDelete

	OwnerId string `json:"owner_id" validate:"required"`
	Name    string `json:"name" validate:"required"`

	Tier *WorkspaceTier `json:"tier" validate:"required"`
}

func (entity *Workspace) TableName() string {
	return "kanthor_workspace"
}

func (entity *Workspace) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("ws")
	}
}

type WorkspaceTier struct {
	WorkspaceId string `json:"workspace_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
}

func (entity *WorkspaceTier) TableName() string {
	return "kanthor_workspace_tier"
}
