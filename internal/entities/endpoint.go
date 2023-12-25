package entities

import (
	"fmt"
	"net/http"

	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
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
	// HTTP: https:://httpbentity.org/post?app=kanthor.webhook
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

func (entity *Endpoint) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", entity.AppId, IdNsApp),
		validator.StringRequired("name", entity.Name),
		validator.StringRequired("secret_key", entity.SecretKey),
		validator.StringLen("secret_key", entity.SecretKey, 16, 32),
		validator.StringUri("uri", entity.Uri),
		validator.StringOneOf("method", entity.Method, []string{http.MethodPost, http.MethodPut}),
	)
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

func (entity *EndpointRule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", entity.EpId, IdNsEp),
		validator.StringRequired("name", entity.Name),
		validator.NumberGreaterThan("priority", entity.Priority, 0),
		validator.StringRequired("condition_source", entity.ConditionSource),
		validator.StringRequired("condition_expression", entity.ConditionExpression),
	)
}
