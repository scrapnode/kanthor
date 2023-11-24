package entities

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/domain/status"
)

var ResponseProps = []string{
	"id",
	"timestamp",
	"msg_id",
	"ep_id",
	"req_id",
	"tier",
	"app_id",
	"type",
	"metadata",
	"headers",
	"body",
	"uri",
	"status",
	"error",
}

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

func (entity *Response) Reschedulable() bool {
	return status.Is5xx(entity.Status) || entity.Status == status.ErrUnknown || entity.Status == status.None
}
