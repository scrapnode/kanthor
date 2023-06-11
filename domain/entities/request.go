package entities

type Request struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Uri      string            `json:"uri"`
	Body     []byte            `json:"body"`
	Metadata map[string]string `json:"metadata"`

	Status int `json:"status"`
}

func (entity *Request) TableName() string {
	return "request"
}
