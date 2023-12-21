package entities

import (
	"encoding/json"
)

type Message struct {
	TSEntity

	Tier     string
	AppId    string
	Type     string
	Metadata Metadata
	Headers  Header
	Body     string
}

func (entity *Message) TableName() string {
	return TableMsg
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

var MessageProps = []string{
	"id",
	"timestamp",
	"tier",
	"app_id",
	"type",
	"metadata",
	"headers",
	"body",
}

var MessageMappers = map[string]func(doc *Message) any{
	"id":        func(doc *Message) any { return doc.Id },
	"timestamp": func(doc *Message) any { return doc.Timestamp },
	"tier":      func(doc *Message) any { return doc.Tier },
	"app_id":    func(doc *Message) any { return doc.AppId },
	"type":      func(doc *Message) any { return doc.Type },
	"metadata":  func(doc *Message) any { return doc.Metadata.String() },
	"headers":   func(doc *Message) any { return doc.Headers.String() },
	"body":      func(doc *Message) any { return doc.Body },
}
