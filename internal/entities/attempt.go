package entities

import (
	"encoding/json"
)

type Attempt struct {
	ReqId string
	MsgId string
	AppId string
	Tier  string

	ScheduledAt int64
	Status      int

	ResId       string
	CompletedAt int64

	ScheduleCounter int
	ScheduleNext    int64
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
	AppId string
	Tier  string
	From  int64
	To    int64
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
	Attemptable map[string]*Attempt
	Ignore      []string
}

var AttemptProps = []string{
	"req_id",
	"msg_id",
	"app_id",
	"tier",
	"status",
	"res_id",
	"schedule_counter",
	"schedule_next",
	"scheduled_at",
	"completed_at",
}
var AttemptMappers = map[string]func(doc *Attempt) any{
	"req_id":           func(doc *Attempt) any { return doc.ReqId },
	"msg_id":           func(doc *Attempt) any { return doc.MsgId },
	"app_id":           func(doc *Attempt) any { return doc.AppId },
	"tier":             func(doc *Attempt) any { return doc.Tier },
	"status":           func(doc *Attempt) any { return doc.Status },
	"res_id":           func(doc *Attempt) any { return doc.ResId },
	"schedule_counter": func(doc *Attempt) any { return doc.ScheduleCounter },
	"schedule_next":    func(doc *Attempt) any { return doc.ScheduleNext },
	"scheduled_at":     func(doc *Attempt) any { return doc.ScheduledAt },
	"completed_at":     func(doc *Attempt) any { return doc.CompletedAt },
}
