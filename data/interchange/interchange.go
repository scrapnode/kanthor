package interchange

import "github.com/scrapnode/kanthor/domain/entities"

type Interchange struct {
	Workspaces []Workspace `json:"workspaces"`
}

type Workspace struct {
	*entities.Workspace
	Applications []Application `json:"applications"`
}

type Application struct {
	*entities.Application
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	*entities.Endpoint
	Rules []EndpointRule `json:"rules"`
}

type EndpointRule struct {
	*entities.EndpointRule
}
