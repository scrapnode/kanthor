package demo

import (
	_ "embed"
	"encoding/json"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func Project(wsId string, bytes []byte) (*ProjectEntities, error) {
	var in struct {
		Projects []project `json:"projects"`
	}
	if err := json.Unmarshal(bytes, &in); err != nil {
		return nil, err
	}

	now := timer.New().Now()
	results := &ProjectEntities{}
	for _, project := range in.Projects {
		for _, application := range project.Applications {
			app := &entities.Application{WorkspaceId: wsId, Name: application.Name}
			app.GenId()
			app.SetAT(now)
			results.Applications = append(results.Applications, *app)

			for _, endpoint := range application.Endpoints {
				ep := &entities.Endpoint{
					AppId:  app.Id,
					Name:   endpoint.Name,
					Method: endpoint.Method,
					Uri:    endpoint.Uri,
				}
				ep.GenId()
				ep.SetAT(now)
				results.Endpoints = append(results.Endpoints, *ep)

				for _, rule := range endpoint.Rules {
					epr := &entities.EndpointRule{
						EndpointId:          ep.Id,
						Name:                rule.Name,
						ConditionSource:     rule.ConditionSource,
						ConditionExpression: rule.ConditionExpression,
					}
					epr.GenId()
					epr.SetAT(now)
					results.EndpointRules = append(results.EndpointRules, *epr)
				}
			}
		}
	}

	return results, nil
}

type ProjectEntities struct {
	Applications  []entities.Application
	Endpoints     []entities.Endpoint
	EndpointRules []entities.EndpointRule
}

type project struct {
	Applications []struct {
		Name string `json:"name" validate:"required"`

		Endpoints []struct {
			Name   string `json:"name" validate:"required"`
			Method string `json:"method" validate:"required,oneof=POST PUT PATCH"`
			Uri    string `json:"uri" validate:"required,uri"`

			Rules []struct {
				Name                string `json:"name" validate:"required"`
				ConditionSource     string `json:"condition_source" validate:"required""`
				ConditionExpression string `json:"condition_expression" validate:"required"`
			} `json:"rules"`
		} `json:"endpoints"`
	} `json:"applications"`
}
