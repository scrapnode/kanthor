package entities

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/infrastructure/utils"
)

var (
	StatusScheduled = 0
)

type Request struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Uri string `json:"uri"`
	// HTTP: POST/PUT/PATCH
	Method   string            `json:"method"`
	Body     []byte            `json:"body"`
	Metadata map[string]string `json:"metadata"`

	Status int `json:"status"`
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
