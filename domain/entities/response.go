package entities

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	RequestId string            `json:"request_id"`
	Metadata  map[string]string `json:"metadata"`

	Uri     string      `json:"uri"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`

	Status int `json:"status"`
}

func (entity *Response) TableName() string {
	return "request"
}

func (entity *Response) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Response) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}
