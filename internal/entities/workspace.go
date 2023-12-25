package entities

import "github.com/scrapnode/kanthor/pkg/validator"

type Workspace struct {
	Entity
	AuditTime

	OwnerId string
	Name    string
	Tier    string
}

func (entity *Workspace) TableName() string {
	return TableWs
}

func (entity *Workspace) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("owner_id", entity.OwnerId),
		validator.StringRequired("name", entity.Name),
		validator.StringRequired("tier", entity.Tier),
	)
}

type WorkspaceCredentials struct {
	Entity
	AuditTime

	WsId      string
	Name      string
	Hash      string
	ExpiredAt int64
}

func (entity *WorkspaceCredentials) TableName() string {
	return TableWsc
}

func (entity *WorkspaceCredentials) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("ws_id", entity.WsId),
		validator.StringRequired("name", entity.Name),
		validator.StringRequired("hash", entity.Hash),
		validator.NumberGreaterThan("expired_at", entity.ExpiredAt, 0),
	)
}
