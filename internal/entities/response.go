package entities

import (
	"encoding/json"
)

type Response struct {
	TSEntity

	EpId  string
	MsgId string
	ReqId string

	Tier     string
	AppId    string
	Type     string
	Metadata Metadata

	Headers Header
	Body    string
	Uri     string
	Status  int
	Error   string
}

func (entity *Response) TableName() string {
	return TableRes
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
