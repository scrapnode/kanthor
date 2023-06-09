package entities

type Workspace struct {
	Entity
	AuditTime
	SoftDelete

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`

	Tier *WorkspaceTier
}

type WorkspaceTier struct {
	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
}
