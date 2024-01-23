package entities

import "encoding/json"

type RecoveryTask struct {
	AppId string
	To    int64
	From  int64

	Init int64
}

func (entity *RecoveryTask) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *RecoveryTask) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *RecoveryTask) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
