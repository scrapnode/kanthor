package rest

import "github.com/scrapnode/kanthor/internal/entities"

type Workspace struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Tier    string `json:"tier"`
} // @name Workspace

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
} // @name WorkspaceCredentials

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

type Message struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`

	AppId    string `json:"app_id"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`
	Headers  string `json:"headers"`
	Body     string `json:"body"`
} // @name Message

func ToMessage(doc *entities.Message) *Message {
	return &Message{
		Id:        doc.Id,
		Timestamp: doc.Timestamp,
		AppId:     doc.AppId,
		Type:      doc.Type,
		Metadata:  doc.Metadata.String(),
		Headers:   doc.Headers.String(),
		Body:      doc.Body,
	}
}
