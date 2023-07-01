package entities

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/http"
)

type Request struct {
	Entity
	TimeSeries
	Tier string `json:"tier"`

	AppId      string `json:"app_id"`
	Type       string `json:"type"`
	EndpointId string `json:"endpoint_id"`

	// HTTP: POST/PUT/PATCH
	Method   string            `json:"method"`
	Uri      string            `json:"uri"`
	Headers  http.Header       `json:"headers"`
	Body     string            `json:"body"`
	Metadata map[string]string `json:"metadata"`
}

func (entity *Request) TableName() string {
	return "request"
}

func (entity *Request) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("req")
	}
}

func (entity *Request) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Request) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}
