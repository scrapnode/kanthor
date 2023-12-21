package entities

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
