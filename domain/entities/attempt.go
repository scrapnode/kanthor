package entities

import (
	"encoding/json"
)

type Attempt struct {
	ReqId string `json:"req_id"`
	MsgId string `json:"msg_id"`
	AppId string `json:"app_id"`
	Tier  string `json:"tier"`

	ScheduledAt int64 `json:"scheduled_at"`
	Status      int   `json:"status"`

	ResId       string `json:"res_id"`
	CompletedAt int64  `json:"completed_at"`

	ScheduleCounter int   `json:"schedule_counter"`
	ScheduleNext    int64 `json:"schedule_next"`
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

func (entity *Attempt) Complete() bool {
	return entity.CompletedAt > 0
}

type AttemptTrigger struct {
	AppId string `json:"app_id"`
	Tier  string `json:"tier"`
	From  int64  `json:"from"`
	To    int64  `json:"to"`
}

func (noti *AttemptTrigger) Marshal() ([]byte, error) {
	return json.Marshal(noti)
}

func (noti *AttemptTrigger) Unmarshal(data []byte) error {
	return json.Unmarshal(data, noti)
}

func (noti *AttemptTrigger) String() string {
	data, _ := json.Marshal(noti)
	return string(data)
}

type AttemptStrive struct {
	Attemptable []Attempt
	Ignore      []string
}
