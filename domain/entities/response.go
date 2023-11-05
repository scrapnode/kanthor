package entities

import (
	"encoding/json"
)

type Response struct {
	TSEntity

	MsgId string `json:"msg_id"`
	EpId  string `json:"ep_id"`
	ReqId string `json:"req_id"`

	Tier     string   `json:"tier"`
	AppId    string   `json:"app_id"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`

	Headers Header `json:"headers"`
	Body    string `json:"body"`
	Uri     string `json:"uri"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func (entity *Response) TableName() string {
	return TableRes
}

func (entity *Response) GenId() {
	if entity.Id == "" {
		entity.Id = ResId()
	}
}

func (entity *Response) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Response) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *Response) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
