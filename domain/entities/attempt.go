package entities

import (
	"encoding/json"
)

type Attempt struct {
	ReqId string `json:"req_id"`

	Tier   string `json:"tier"`
	Status int    `json:"status"`
	ResId  string `json:"res_id"`

	ScheduleCounter int   `json:"schedule_counter"`
	ScheduleNext    int64 `json:"schedule_next"`
	ScheduledAt     int64 `json:"scheduled_at"`
	CompletedAt     int64 `json:"completed_at"`
}

func (entity *Attempt) TableName() string {
	return "kanthor_attempt"
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
