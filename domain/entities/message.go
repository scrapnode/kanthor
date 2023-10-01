package entities

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/pkg/utils"
)

// Messages are allocated based on bucket
// For SQL: create a composite index for AppId+Type+Bucket, sort by ID (ksuid)
// for Dynamo-style: partition by AppId+Type+Bucket, sort by ID (ksuid)
type Message struct {
	TSEntity

	Tier     string   `json:"tier"`
	AppId    string   `json:"app_id"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`
	Headers  Header   `json:"headers"`
	Body     []byte   `json:"body"`
}

func (entity *Message) TableName() string {
	return "kanthor_message"
}

func (entity *Message) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("msg")
	}
}

func (entity *Message) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Message) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *Message) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
