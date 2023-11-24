package entities

import (
	"encoding/json"
)

var RequestProps = []string{
	"id",
	"timestamp",
	"msg_id",
	"ep_id",
	"tier",
	"app_id",
	"type",
	"metadata",
	"headers",
	"body",
	"uri",
	"method",
}

type Request struct {
	TSEntity

	MsgId string `json:"msg_id"`
	EpId  string `json:"ep_id"`

	Tier     string   `json:"tier"`
	AppId    string   `json:"app_id"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`

	// HTTP: POST/PUT/PATCH
	Headers Header `json:"headers"`
	Body    string `json:"body"`
	Uri     string `json:"uri"`
	Method  string `json:"method"`
}

func (entity *Request) TableName() string {
	return TableReq
}

func (entity *Request) GenId() {
	if entity.Id == "" {
		entity.Id = ReqId()
	}
}

func (entity *Request) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Request) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *Request) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
