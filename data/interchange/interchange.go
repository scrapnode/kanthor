package interchange

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/scrapnode/kanthor/pkg/validator"
)

func Unmarshal(data []byte) (*Workspace, error) {
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	if err := ws.Validate(); err != nil {
		return nil, err
	}

	return &ws, nil
}

type Workspace struct {
	Applications []Application `json:"applications"`
}

func (ws *Workspace) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("data.interchange.applications", ws.Applications),
	)
	if err != nil {
		return err
	}
	for i, app := range ws.Applications {
		if err := app.Validate(fmt.Sprintf("applications[%d]", i)); err != nil {
			return err
		}
	}
	return nil
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Application struct {
	Name      string     `json:"name"`
	Endpoints []Endpoint `json:"endpoints"`
}

func (app *Application) Validate(key string) error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired(fmt.Sprintf("%s.name", key), app.Name),
		validator.SliceRequired(fmt.Sprintf("%s.endpoints", key), app.Endpoints),
	)
	if err != nil {
		return err
	}
	for i, ep := range app.Endpoints {
		if err := ep.Validate(fmt.Sprintf("%s.endpoints[%d]", key, i)); err != nil {
			return err
		}
	}
	return nil
}

type Endpoint struct {
	Name string `json:"name"`

	SecretKey string         `json:"secret_key"`
	Method    string         `json:"method"`
	Uri       string         `json:"uri"`
	Rules     []EndpointRule `json:"rules"`
}

func (ep *Endpoint) Validate(key string) error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired(fmt.Sprintf("%s.name", key), ep.Name),
		validator.StringLenIfNotEmpty(fmt.Sprintf("%s.secret_key", key), ep.SecretKey, 16, 32),
		validator.StringUri(fmt.Sprintf("%s.uri", key), ep.Uri),
		validator.StringOneOf(fmt.Sprintf("%s.method", key), ep.Method, []string{http.MethodPost, http.MethodPut}),
	)
	if err != nil {
		return err
	}
	for i, epr := range ep.Rules {
		if err := epr.Validate(fmt.Sprintf("%s.rules[%d]", key, i)); err != nil {
			return err
		}
	}
	return nil
}

type EndpointRule struct {
	Name                string `json:"name"`
	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source"`
	ConditionExpression string `json:"condition_expression"`
}

func (epr *EndpointRule) Validate(key string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired(fmt.Sprintf("%s.name", key), epr.Name),
		validator.NumberGreaterThanOrEqual(fmt.Sprintf("%s.priority", key), epr.Priority, 0),
		validator.StringRequired(fmt.Sprintf("%s.condition_source", key), epr.ConditionSource),
		validator.StringRequired(fmt.Sprintf("%s.condition_expression", key), epr.ConditionExpression),
	)
}
