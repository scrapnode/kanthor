package entities

import (
	"github.com/scrapnode/kanthor/infrastructure/utils"
	"net/http"
)

type Response struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	RequestId string `json:"request_id"`

	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`

	Uri    string `json:"uri"`
	Status int    `json:"status"`
}

func (entity *Response) GenId() {
	entity.Id = utils.ID("res")
}
