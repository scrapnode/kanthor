package entities

import (
	"github.com/scrapnode/kanthor/infrastructure/utils"
)

type Response struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	RequestId string `json:"request_id"`

	Uri      string            `json:"uri"`
	Metadata map[string]string `json:"metadata"`
	Body     []byte            `json:"body"`

	Status int `json:"status"`
}

func (entity *Response) GenId() {
	entity.Id = utils.ID("res")
}
