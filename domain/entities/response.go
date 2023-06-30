package entities

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/http"
)

type Response struct {
	Entity
	TimeSeries
	Tier string `json:"tier"`

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Metadata map[string]string `json:"metadata"`

	Uri     string      `json:"uri"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`

	Status int `json:"status"`
}

func (entity *Response) TableName() string {
	return "request"
}

func (entity *Response) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("res")
	}
}

func (entity *Response) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Response) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}
