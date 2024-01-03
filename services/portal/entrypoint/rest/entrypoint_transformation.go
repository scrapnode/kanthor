package rest

import (
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
)

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

	Headers string `json:"headers"`
	Body    string `json:"body"`
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

	Requests  []Request  `json:"requests"`
	Responses []Response `json:"responses"`
} // @name EndpointMessage

func ToEndpointMessage(doc *entities.EndpointMessage, requests []entities.Request, responses []entities.Response) *EndpointMessage {
	msg := &EndpointMessage{
		Message:          ToMessage(&doc.Message),
		RequestCount:     doc.RequestCount,
		RequestLatestTs:  doc.RequestLatestTs,
		ResponseCount:    doc.ResponseCount,
		ResponseLatestTs: doc.ResponseLatestTs,
		SuccessId:        doc.SuccessId,

		Requests:  make([]Request, 0),
		Responses: make([]Response, 0),
	}

	if len(requests) > 0 {
		for _, request := range requests {
			msg.Requests = append(msg.Requests, Request{
				Id:        request.Id,
				Timestamp: request.Timestamp,
				EpId:      request.EpId,
				MsgId:     request.MsgId,
				AppId:     request.AppId,
				Type:      request.Type,
				Metadata:  request.Metadata.String(),
				Headers:   request.Headers.String(),
				Body:      request.Body,
				Uri:       request.Uri,
				Method:    request.Method,
			})
		}
	}

	if len(responses) > 0 {
		for _, response := range responses {
			msg.Responses = append(msg.Responses, Response{
				Id:        response.Id,
				Timestamp: response.Timestamp,
				EpId:      response.EpId,
				MsgId:     response.MsgId,
				AppId:     response.AppId,
				Type:      response.Type,
				Metadata:  response.Metadata.String(),
				Headers:   response.Headers.String(),
				Body:      response.Body,
				Uri:       response.Uri,
				Status:    response.Status,
				Error:     response.Error,
			})
		}
	}

	return msg
}

type Request struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	EpId      string `json:"ep_id"`
	MsgId     string `json:"msg_id"`

	AppId    string `json:"app_id"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`

	Headers string `json:"headers"`
	Body    string `json:"body"`
	Uri     string `json:"uri"`
	Method  string `json:"method"`
} // @name Request

type Response struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	EpId      string `json:"ep_id"`
	MsgId     string `json:"msg_id"`
	ReqId     string `json:"req_id"`

	AppId    string `json:"app_id"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`

	Headers string `json:"headers"`
	Body    string `json:"body"`
	Uri     string `json:"uri"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
} // @name Response

type WorkspaceSnapshot struct {
	Name         string                 `json:"name"`
	Applications []WorkspaceSnapshotApp `json:"applications"`
} // @name WorkspaceSnapshot

type WorkspaceSnapshotApp struct {
	Name      string                `json:"name"`
	Endpoints []WorkspaceSnapshotEp `json:"endpoints"`
} // @name WorkspaceSnapshotApp

type WorkspaceSnapshotEp struct {
	Name   string                 `json:"name"`
	Method string                 `json:"method"`
	Uri    string                 `json:"uri"`
	Rules  []WorkspaceSnapshotEpr `json:"rules"`
} // @name WorkspaceSnapshotEp

type WorkspaceSnapshotEpr struct {
	Name                string `json:"name"`
	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source"`
	ConditionExpression string `json:"condition_expression"`
} // @name WorkspaceSnapshotEpr

func ToWorkspaceSnapshot(snapshot *entities.WorkspaceSnapshot) *WorkspaceSnapshot {
	returning := &WorkspaceSnapshot{Name: snapshot.Name}

	for _, app := range snapshot.Applications {
		application := WorkspaceSnapshotApp{Name: app.Name}
		for _, ep := range app.Endpoints {
			endpoint := WorkspaceSnapshotEp{
				Name:   ep.Name,
				Method: ep.Method,
				Uri:    ep.Uri,
			}
			for _, epr := range ep.Rules {
				rule := WorkspaceSnapshotEpr{
					Name:                epr.Name,
					Priority:            epr.Priority,
					Exclusionary:        epr.Exclusionary,
					ConditionSource:     epr.ConditionSource,
					ConditionExpression: epr.ConditionExpression,
				}
				endpoint.Rules = append(endpoint.Rules, rule)
			}
			application.Endpoints = append(application.Endpoints, endpoint)
		}
		returning.Applications = append(returning.Applications, application)
	}

	return returning
}

func FromWorkspaceSnapshot(snapshot *WorkspaceSnapshot, id string) *entities.WorkspaceSnapshot {
	returning := &entities.WorkspaceSnapshot{
		Id:           id,
		Name:         snapshot.Name,
		Applications: make(map[string]entities.WorkspaceSnapshotApp),
	}

	for _, app := range snapshot.Applications {
		application := entities.WorkspaceSnapshotApp{
			Name:      app.Name,
			Endpoints: make(map[string]entities.WorkspaceSnapshotEp),
		}

		for _, ep := range app.Endpoints {
			endpoint := entities.WorkspaceSnapshotEp{
				Name:   ep.Name,
				Method: ep.Method,
				Uri:    ep.Uri,
				Rules:  make(map[string]entities.WorkspaceSnapshotEpr),
			}

			for _, epr := range ep.Rules {
				rule := entities.WorkspaceSnapshotEpr{
					Name:                epr.Name,
					Priority:            epr.Priority,
					Exclusionary:        epr.Exclusionary,
					ConditionSource:     epr.ConditionSource,
					ConditionExpression: epr.ConditionExpression,
				}
				endpoint.Rules[suid.New(entities.IdNsEpr)] = rule
			}

			application.Endpoints[suid.New(entities.IdNsEp)] = endpoint
		}

		returning.Applications[suid.New(entities.IdNsApp)] = application
	}

	return returning
}
