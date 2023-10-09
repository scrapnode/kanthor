package planner

import (
	"github.com/scrapnode/kanthor/domain/entities"
)

var (
	ConditionSourceAppId = "app_id"
	ConditionSourceType  = "type"
)

func ConditionSource(rule *entities.EndpointRule, msg *entities.Message) string {
	if rule.ConditionSource == ConditionSourceAppId {
		return msg.AppId
	}
	if rule.ConditionSource == ConditionSourceType {
		return msg.Type
	}
	return ""
}
