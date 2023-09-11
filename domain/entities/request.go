package entities

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/pkg/utils"
)

type Request struct {
	TSEntity
	AttId string `json:"attempt_id"`

	Tier     string   `json:"tier"`
	AppId    string   `json:"app_id"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`

	// HTTP: POST/PUT/PATCH
	Headers Header `json:"headers"`
	Body    []byte `json:"body"`
	Uri     string `json:"uri"`
	Method  string `json:"method"`
}

func (entity *Request) TableName() string {
	return "kanthor_request"
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
