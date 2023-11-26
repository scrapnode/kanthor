package entities

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/internal/domain/status"
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

func (entity *Response) Reschedulable() bool {
	return status.Is5xx(entity.Status) || entity.Status == status.ErrUnknown || entity.Status == status.None
}

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

var ResponseMappers = map[string]func(doc *Response) any{
	"id":        func(doc *Response) any { return doc.Id },
	"timestamp": func(doc *Response) any { return doc.Timestamp },
	"msg_id":    func(doc *Response) any { return doc.MsgId },
	"ep_id":     func(doc *Response) any { return doc.EpId },
	"req_id":    func(doc *Response) any { return doc.ReqId },
	"tier":      func(doc *Response) any { return doc.Tier },
	"app_id":    func(doc *Response) any { return doc.AppId },
	"type":      func(doc *Response) any { return doc.Type },
	"metadata":  func(doc *Response) any { return doc.Metadata.String() },
	"headers":   func(doc *Response) any { return doc.Headers.String() },
	"body":      func(doc *Response) any { return doc.Body },
	"uri":       func(doc *Response) any { return doc.Uri },
	"status":    func(doc *Response) any { return doc.Status },
	"error":     func(doc *Response) any { return doc.Error },
}
