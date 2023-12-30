package rest

import "github.com/scrapnode/kanthor/internal/entities"

type Application struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	WsId      string `json:"ws_id"`
	Name      string `json:"name"`
} // @name Application

func ToApplication(doc *entities.Application) *Application {
	return &Application{
		Id:        doc.Id,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		WsId:      doc.WsId,
		Name:      doc.Name,
	}
}

type Endpoint struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	AppId     string `json:"app_id"`
	Name      string `json:"name"`
	Method    string `json:"method"`
	Uri       string `json:"uri"`
} // @name Endpoint

func ToEndpoint(doc *entities.Endpoint) *Endpoint {
	return &Endpoint{
		Id:        doc.Id,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		AppId:     doc.AppId,
		Name:      doc.Name,
		Method:    doc.Method,
		Uri:       doc.Uri,
	}
}

type EndpointRule struct {
	Id                  string `json:"id"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
	EpId                string `json:"ep_id"`
	Name                string `json:"name"`
	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source"`
	ConditionExpression string `json:"condition_expression"`
} // @name EndpointRule

func ToEndpointRule(doc *entities.EndpointRule) *EndpointRule {
	return &EndpointRule{
		Id:                  doc.Id,
		CreatedAt:           doc.CreatedAt,
		UpdatedAt:           doc.UpdatedAt,
		EpId:                doc.EpId,
		Name:                doc.Name,
		Priority:            doc.Priority,
		Exclusionary:        doc.Exclusionary,
		ConditionSource:     doc.ConditionSource,
		ConditionExpression: doc.ConditionExpression,
	}
}
