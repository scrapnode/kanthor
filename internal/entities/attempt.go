package entities

import "encoding/json"

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

type Attempt struct {
	ReqId string

	MsgId string
	AppId string
	Tier  string

	Status int

	SuccessId       string
	ScheduleCounter int
	ScheduleNext    int64
	ScheduledAt     int64
	CompletedAt     int64
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
