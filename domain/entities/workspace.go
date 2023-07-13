package entities

import "github.com/scrapnode/kanthor/pkg/utils"

type Workspace struct {
	Entity
	AuditTime
	SoftDelete

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`

	Tier *WorkspaceTier
}

func (entity *Workspace) TableName() string {
	return "workspace"
}

func (entity *Workspace) GenId() {
	entity.Id = utils.ID("ws")
}

type WorkspaceTier struct {
	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
}

func (entity *WorkspaceTier) TableName() string {
	return "workspace_tier"
}

type WorkspacePrivilege struct {
	Entity
	AuditTime
	SoftDelete

	WorkspaceId string `json:"workspace_id"`
	AccountSub  string `json:"account_sub"`
	AccountRole string `json:"account_role"`
}

func (entity *WorkspacePrivilege) TableName() string {
	return "workspace_privilege"
}

func (entity *WorkspacePrivilege) GenId() {
	entity.Id = utils.ID("wsp")
}
