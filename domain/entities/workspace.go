package entities

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

type WorkspaceTier struct {
	WorkspaceId string `json:"workspace_id"`
	Name        string `json:"name"`
}

func (entity *WorkspaceTier) TableName() string {
	return "workspace_tier"
}

func DefaultTier(id string) *WorkspaceTier {
	return &WorkspaceTier{WorkspaceId: id, Name: "default"}
}
