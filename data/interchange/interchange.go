package interchange

import "github.com/scrapnode/kanthor/domain/entities"

type Interchange struct {
	Workspaces []Workspace `json:"workspaces" validate:"required,dive,required"`
}

type Workspace struct {
	*entities.Workspace
	Tier         *entities.WorkspaceTier
	Credentials  []entities.WorkspaceCredentials `json:"credentials" validate:"required"`
	Applications []Application                   `json:"applications" validate:"required"`
}

type Application struct {
	*entities.Application
	Endpoints []Endpoint `json:"endpoints" validate:"required"`
}

type Endpoint struct {
	*entities.Endpoint
	Rules []EndpointRule `json:"rules" validate:"required"`
}

type EndpointRule struct {
	*entities.EndpointRule
}
