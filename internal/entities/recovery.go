package entities

import "encoding/json"

type Recovery struct {
	AppId string
	To    int64
	From  int64

	Init int64
}

func (entity *Recovery) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Recovery) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *Recovery) String() string {
	data, _ := json.Marshal(entity)
	return string(data)
}
