package entities

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/pkg/utils"
)

// Requests are allocated based on bucket
// For SQL: create a composite index for AppId+Type+Bucket, sort by ID (ksuid)
// for Dynamo-style: partition by AppId+Type+Bucket, sort by ID (ksuid)
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

func (entity *Request) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
