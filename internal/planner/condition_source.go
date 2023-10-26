package planner

import (
	"github.com/scrapnode/kanthor/domain/entities"
)

var (
	ConditionSourceType = "type"
)

func ConditionSource(rule *entities.EndpointRule, msg *entities.Message) string {
	if rule.ConditionSource == ConditionSourceType {
		return msg.Type
	}
	return ""
}
