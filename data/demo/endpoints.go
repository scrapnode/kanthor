package demo

import (
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
)

func Endpoints(appId string) []entities.Endpoint {
	return []entities.Endpoint{
		{
			AppId:  appId,
			Name:   "POST httpbin.org/post",
			Method: "POST",
			Uri:    "https:://httpbin.org/post?app=kanthor",
		},
		{
			AppId:  appId,
			Name:   "PATCH httpbin.org/status/500",
			Method: "PATCH",
			Uri:    "https:://httpbin.org/status/500?app=kanthor",
		},
	}
}

func EndpointRules(appId string, epIds []string) []entities.EndpointRule {
	var rules []entities.EndpointRule

	for _, epId := range epIds {
		rules = append(rules, entities.EndpointRule{
			EndpointId:          epId,
			Name:                "match all request to app",
			ConditionSource:     "app_id",
			ConditionExpression: fmt.Sprintf("equal::%s", appId),
		})
	}

	return rules
}
