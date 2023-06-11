package entities

import "github.com/scrapnode/kanthor/infrastructure/utils"

type Endpoint struct {
	Entity
	AuditTime
	SoftDelete

	AppId string `json:"app_id"`

	Name string `json:"name"`
	// format: scheme ":" ["//" authority] path ["?" query] ["#" fragment]
	// HTTP: https:://httpbin.org/post?app=kanthor.webhook
	// gRPC: grpc:://app.kanthorlabs.com
	Uri string `json:"uri"`
	// HTTP: POST/PUT/PATCH
	Method string `json:"method"`
}

func (entity *Endpoint) TableName() string {
	return "endpoint"
}

func (entity *Endpoint) GenId() {
	entity.Id = utils.ID("ep")
}

type EndpointRule struct {
	Entity
	AuditTime
	SoftDelete

	EndpointId string `json:"endpoint_id"`

	// examples:
	// - regex::.*
	// - type::orders.paid
	ConditionSource     string `json:"condition_source"`
	ConditionExpression string `json:"condition_expression"`
	Priority            int    `json:"priority"`
	// the logic of not-false is true should be used here
	// to guarantee default all rule will be on include mode
	Exclusionary bool `json:"exclusionary"`
}

func (entity *EndpointRule) TableName() string {
	return "endpoint_rule"
}
