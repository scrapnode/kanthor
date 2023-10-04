package entities

import (
	"encoding/json"
)

type Attempt struct {
	Id    string `json:"id"`
	EpId  string `json:"ep_id"`
	ResId string `json:"res_id"`

	Tier   string `json:"tier"`
	Status int    `json:"status"`

	ScheduledAt int64 `json:"scheduled_at"`
	CompletedAt int64 `json:"completed_at"`
}

func (entity *Attempt) TableName() string {
	return "kanthor_attempt"
}

func (entity *Attempt) GenId() {
	if entity.Id == "" {
		entity.Id = AttId()
	}
}

func (entity *Attempt) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Attempt) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *Attempt) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
