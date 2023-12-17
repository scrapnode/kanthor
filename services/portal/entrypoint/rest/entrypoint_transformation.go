package rest

import "github.com/scrapnode/kanthor/internal/entities"

type Workspace struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Tier    string `json:"tier"`
}

func ToWorkspace(doc *entities.Workspace) *Workspace {
	return &Workspace{
		Id:        doc.Id,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		OwnerId:   doc.OwnerId,
		Name:      doc.Name,
		Tier:      doc.Tier,
	}
}

type WorkspaceCredentials struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	WsId      string `json:"ws_id"`
	Name      string `json:"name"`
	ExpiredAt int64  `json:"expired_at"`
}

func ToWorkspaceCredentials(doc *entities.WorkspaceCredentials) *WorkspaceCredentials {
	return &WorkspaceCredentials{
		Id:        doc.Id,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		WsId:      doc.WsId,
		Name:      doc.Name,
		ExpiredAt: doc.ExpiredAt,
	}
}
