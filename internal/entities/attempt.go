package entities

import (
	"encoding/json"
	"fmt"
)

type AttemptTask struct {
	AppId string
	EpId  string
	To    int64
	From  int64

	Init int64
}

func (entity *AttemptTask) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *AttemptTask) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *AttemptTask) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}

type AttemptTrigger struct {
	To   int64
	From int64

	Init int64
}

func (entity *AttemptTrigger) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *AttemptTrigger) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *AttemptTrigger) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}

type Attempt struct {
	ReqId string

	MsgId string
	EpId  string
	AppId string
	Tier  string

	AttemptState
}

func (entity *Attempt) Id() string {
	return fmt.Sprintf("%s/%s/%s/%d", entity.MsgId, entity.EpId, entity.ReqId, entity.ScheduleCounter)
}

func (entity *Attempt) TableName() string {
	return TableAtt
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

type AttemptState struct {
	ScheduleCounter int    `json:"schedule_counter"`
	ScheduleNext    int64  `json:"schedule_next"`
	ScheduledAt     int64  `json:"scheduled_at"`
	CompletedAt     int64  `json:"completed_at"`
	CompletedId     string `json:"completed_id"`
}

func (entity *AttemptState) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
