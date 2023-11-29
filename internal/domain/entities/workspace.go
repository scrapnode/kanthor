package entities

type Workspace struct {
	Entity
	AuditTime

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Tier    string `json:"tier"`
}

func (entity *Workspace) TableName() string {
	return TableWs
}

type WorkspaceCredentials struct {
	Entity
	AuditTime

	WsId      string `json:"ws_id"`
	Name      string `json:"name"`
	Hash      string `json:"hash,omitempty"`
	ExpiredAt int64  `json:"expired_at"`
}

func (entity *WorkspaceCredentials) TableName() string {
	return TableWsc
}
