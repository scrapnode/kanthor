package entities

import (
	"github.com/scrapnode/kanthor/pkg/utils"
)

type Endpoint struct {
	Entity
	AuditTime

	AppId string `json:"app_id" validate:"required"`
	Name  string `json:"name" validate:"required"`

	// HTTP: POST/PUT/PATCH
	Method string `json:"method" validate:"required,oneof=POST PUT PATCH"`
	// format: scheme ":" ["//" authority] path ["?" query] ["#" fragment]
	// HTTP: https:://httpbin.org/post?app=kanthor.webhook
	// gRPC: grpc:://app.kanthorlabs.com
	Uri string `json:"uri" validate:"required,uri"`
}

func (entity *Endpoint) TableName() string {
	return "kanthor_endpoint"
}

func (entity *Endpoint) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("ep")
	}
}

type EndpointRule struct {
	Entity
	AuditTime

	EndpointId string `json:"endpoint_id" validate:"required"`
	Name       string `json:"name" validate:"required"`

	Priority int `json:"priority"`
	// the logic of not-false is true should be used here
	// to guarantee default all rule will be on include mode
	Exclusionary bool `json:"exclusionary"`

	// examples
	//  - app_id
	//  - type
	//  - body
	//  - metadata
	ConditionSource string `json:"condition_source" validate:"required"`
	// examples:
	// 	- equal::orders.paid
	// 	- regex::.*
	ConditionExpression string `json:"condition_expression" validate:"required"`
}

func (entity *EndpointRule) TableName() string {
	return "kanthor_endpoint_rule"
}

func (entity *EndpointRule) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("epr")
	}
}
