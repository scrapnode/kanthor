package entities

import "encoding/json"

type Response struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	RequestId string `json:"request_id"`

	Uri      string            `json:"uri"`
	Metadata map[string]string `json:"metadata"`
	Body     []byte            `json:"body"`

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
