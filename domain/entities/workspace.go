package entities

type Workspace struct {
	Entity
	AuditTime

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Tier    string `json:"tier"`
}

func (entity *Workspace) TableName() string {
	return "kanthor_workspace"
}

func (entity *Workspace) GenId() {
	if entity.Id == "" {
		entity.Id = WsId()
	}
}

type WorkspaceCredentials struct {
	Entity
	AuditTime

	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
	Hash        string `json:"hash,omitempty"`
	ExpiredAt   int64  `json:"expired_at"`
}

func (entity *WorkspaceCredentials) TableName() string {
	return "kanthor_workspace_credentials"
}

func (entity *WorkspaceCredentials) GenId() {
	if entity.Id == "" {
		entity.Id = WscId()
	}
}
