package entities

import (
	"github.com/scrapnode/kanthor/pkg/utils"
)

type Workspace struct {
	Entity
	AuditTime

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
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
	Entity
	AuditTime

	WorkspaceId string `json:"workspace_id" validate:"required,startswith=ws_"`
	Name        string `json:"name" validate:"required"`
}

func (entity *WorkspaceTier) TableName() string {
	return "kanthor_workspace_tier"
}

func (entity *WorkspaceTier) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("wst")
	}
}

type WorkspaceCredentials struct {
	Entity
	AuditTime

	WorkspaceId string `json:"workspace_id" validate:"required,startswith=ws_"`
	Name        string `json:"name" validate:"required"`
	Hash        string `json:"hash" validate:"required"`
	ExpiredAt   int64  `json:"expired_at" validate:"required,gt=0"`
}

func (entity *WorkspaceCredentials) TableName() string {
	return "kanthor_workspace_credentials"
}

func (entity *WorkspaceCredentials) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("wsc")
	}
}
