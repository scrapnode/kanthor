package entities

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/pkg/utils"
)

// Responses are allocated based on bucket
// For SQL: create a composite index for AppId+Type+Bucket, sort by ID (ksuid)
// for Dynamo-style: partition by AppId+Type+Bucket, sort by ID (ksuid)
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
	Body    []byte `json:"body"`
	Uri     string `json:"uri"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func (entity *Response) TableName() string {
	return "kanthor_response"
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
