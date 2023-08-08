package entities

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/http"
)

// Message are allocated based on bucket
// For SQL: create a composite index for AppId+Type+Bucket, sort by ID (ksuid)
// for Dynamo-style: partition by AppId+Type+Bucket, sort by ID (ksuid)
// we don't need workspace_id because most time we only retrieve message of app, not of workspace
type Message struct {
	TSEntity

	Tier     string   `json:"tier"`
	AppId    string   `json:"app_id"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`

	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
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
