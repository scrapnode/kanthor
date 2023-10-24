package entities

import (
	"encoding/json"
)

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
		entity.Id = MsgId()
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
