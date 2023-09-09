package interchange

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/infrastructure/validator"
)

func Unmarshal(data []byte) (*Workspace, error) {
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(ws); err != nil {
		return nil, err
	}

	return &ws, nil
}

type Workspace struct {
	Applications []Application `json:"applications" validate:"required,gt=0"`
}

type Application struct {
	Name      string     `json:"name" validate:"required"`
	Endpoints []Endpoint `json:"endpoints" validate:"required,gt=0"`
}

type Endpoint struct {
	Name string `json:"name" validate:"required"`

	SecretKey string         `json:"secret_key"`
	Method    string         `json:"method" validate:"required,oneof=GET POST PUT PATCH"`
	Uri       string         `json:"uri" validate:"required,uri"`
	Rules     []EndpointRule `json:"rules" validate:"required,gt=0"`
}

type EndpointRule struct {
	Name                string `json:"name" validate:"required"`
	Priority            int32  `json:"priority" validate:"required,gte=0"`
	Exclusionary        bool   `json:"exclusionary" validate:"boolean"`
	ConditionSource     string `json:"condition_source" validate:"required"`
	ConditionExpression string `json:"condition_expression" validate:"required"`
}
