package entities

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/utils"
)

type Endpoint struct {
	Entity
	AuditTime
	// @TODO: add deactivated_at
	// DeactivatedAt int64

	AppId     string
	Name      string
	SecretKey string
	// HTTP: POST/PUT/PATCH
	Method string
	// format: scheme ":" ["//" authority] path ["?" query] ["#" fragment]
	// HTTP: https:://httpbin.org/post?app=kanthor.webhook
	// gRPC: grpc:://app.kanthorlabs.com
	Uri string
}

func (entity *Endpoint) TableName() string {
	return TableEp
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
	// DeactivatedAt int64

	EpId string
	Name string

	Priority int32
	// the logic of not-false is true should be used here
	// to guarantee default all rule will be on include mode
	Exclusionary bool

	// examples
	//  - app_id
	//  - type
	//  - body
	//  - metadata
	ConditionSource string
	// examples:
	// 	- equal::orders.paid
	// 	- regex::.*
	ConditionExpression string
}

func (entity *EndpointRule) TableName() string {
	return TableEpr
}
