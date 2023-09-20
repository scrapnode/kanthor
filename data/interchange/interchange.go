package interchange

import (
	"encoding/json"
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
	err := validator.Validate(validator.DefaultConfig, validator.SliceRequired("applications", ws.Applications))
	if err != nil {
		return err
	}
	for _, app := range ws.Applications {
		if err := app.Validate(); err != nil {
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

func (app *Application) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("name", app.Name),
		validator.SliceRequired("endpoints", app.Endpoints),
	)
	if err != nil {
		return err
	}
	for _, ep := range app.Endpoints {
		if err := ep.Validate(); err != nil {
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

func (ep *Endpoint) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("name", ep.Name),
		validator.StringRequired("secret_key", ep.SecretKey),
		validator.StringLen("secret_key", ep.SecretKey, 16, 32),
		validator.StringUri("uri", ep.Uri),
		validator.StringOneOf("method", ep.Method, []string{http.MethodPost, http.MethodPut}),
	)
	if err != nil {
		return err
	}
	for _, epr := range ep.Rules {
		if err := epr.Validate(); err != nil {
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

func (epr *EndpointRule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("name", epr.Name),
		validator.NumberGreaterThan("priority", epr.Priority, 0),
		validator.StringRequired("condition_source", epr.ConditionSource),
		validator.StringRequired("condition_expression", epr.ConditionExpression),
	)
}
