package entities

import (
	"encoding/json"
)

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

var RequestMappers = map[string]func(doc *Request) any{
	"id":        func(doc *Request) any { return doc.Id },
	"timestamp": func(doc *Request) any { return doc.Timestamp },
	"msg_id":    func(doc *Request) any { return doc.MsgId },
	"ep_id":     func(doc *Request) any { return doc.EpId },
	"tier":      func(doc *Request) any { return doc.Tier },
	"app_id":    func(doc *Request) any { return doc.AppId },
	"type":      func(doc *Request) any { return doc.Type },
	"metadata":  func(doc *Request) any { return doc.Metadata.String() },
	"headers":   func(doc *Request) any { return doc.Headers.String() },
	"body":      func(doc *Request) any { return doc.Body },
	"uri":       func(doc *Request) any { return doc.Uri },
	"method":    func(doc *Request) any { return doc.Method },
}
