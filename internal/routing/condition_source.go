package routing

import (
	"github.com/scrapnode/kanthor/internal/entities"
)

var (
	ConditionSourceType  = "type"
	ConditionSourceAppId = "app_id"
)

func ConditionSource(rule *entities.EndpointRule, msg *entities.Message) string {
	if rule.ConditionSource == ConditionSourceType {
		return msg.Type
	}
	if rule.ConditionSource == ConditionSourceAppId {
		return msg.AppId
	}
	return ""
}
