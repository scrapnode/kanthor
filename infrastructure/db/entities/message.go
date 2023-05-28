package entities

import (
	"github.com/scrapnode/kanthor/infrastructure/utils"
	"net/http"
)

// Message are allocated based on bucket
// For SQL: create a composite index for Bucket+AppId+Type, sort by ID (ksuid)
// for Dynamo-style: partition by Bucket+AppId+Type, sort by ID (ksuid)
// we don't need workspace_id because most time we only retrieve message of app, not of workspace
type Message struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

func (entity *Message) GenId() {
	entity.Id = utils.ID("msg")
}
