package entities

import (
	"github.com/scrapnode/kanthor/infrastructure/utils"
	"net/http"
)

type Request struct {
	Entity
	TimeSeries

	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`

	Uri    string `json:"uri"`
	Status int    `json:"status"`
}

func (entity *Request) GenId() {
	entity.Id = utils.ID("req")
}
