package entities

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/utils"
)

type Endpoint struct {
	Entity
	AuditTime
	// @TODO: add deactivated_at
	// DeactivatedAt int64 `json:"deactivated_at"`

	AppId     string `json:"app_id"`
	Name      string `json:"name"`
	SecretKey string `json:"secret_key"`
	// HTTP: POST/PUT/PATCH
	Method string `json:"method"`
	// format: scheme ":" ["//" authority] path ["?" query] ["#" fragment]
	// HTTP: https:://httpbin.org/post?app=kanthor.webhook
	// gRPC: grpc:://app.kanthorlabs.com
	Uri string `json:"uri"`
}

func (entity *Endpoint) TableName() string {
	return "kanthor_endpoint"
}

func (entity *Endpoint) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("ep")
	}
}

func (entity *Endpoint) GenSecretKey() {
	if entity.SecretKey == "" {
		entity.SecretKey = fmt.Sprintf("epsk_%s", utils.RandomString(32))
	}
}

type EndpointRule struct {
	Entity
	AuditTime
	// @TODO: add deactivated_at
	// DeactivatedAt int64 `json:"deactivated_at"`

	EndpointId string `json:"endpoint_id"`
	Name       string `json:"name"`

	Priority int32 `json:"priority"`
	// the logic of not-false is true should be used here
	// to guarantee default all rule will be on include mode
	Exclusionary bool `json:"exclusionary"`

	// examples
	//  - app_id
	//  - type
	//  - body
	//  - metadata
	ConditionSource string `json:"condition_source"`
	// examples:
	// 	- equal::orders.paid
	// 	- regex::.*
	ConditionExpression string `json:"condition_expression"`
}

func (entity *EndpointRule) TableName() string {
	return "kanthor_endpoint_rule"
}

func (entity *EndpointRule) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("epr")
	}
}
