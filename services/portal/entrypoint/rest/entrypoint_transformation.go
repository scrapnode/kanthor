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

type EndpointMessage struct {
	*Message

	RequestCount     int    `json:"request_count"`
	RequestLatestTs  int64  `json:"request_latest_ts"`
	ResponseCount    int    `json:"response_count"`
	ResponseLatestTs int64  `json:"response_latest_ts"`
	SuccessId        string `json:"success_id"`
} // @name EndpointMessage

func ToEndpointMessage(doc *entities.EndpointMessage) *EndpointMessage {
	return &EndpointMessage{
		Message:          ToMessage(&doc.Message),
		RequestCount:     doc.RequestCount,
		RequestLatestTs:  doc.RequestLatestTs,
		ResponseCount:    doc.ResponseCount,
		ResponseLatestTs: doc.ResponseLatestTs,
		SuccessId:        doc.SuccessId,
	}
}
