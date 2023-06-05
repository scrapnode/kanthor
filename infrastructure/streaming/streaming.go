package streaming

import "github.com/scrapnode/kanthor/infrastructure/utils"

var (
	MetaAppId = "KANTHOR_META_APP_ID"
	MetaType  = "KANTHOR_META_TYPE"
	MetaId    = "KANTHOR_META_ID"
)

type Event struct {
	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Id       string            `json:"id"`
	Data     []byte            `json:"data"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Event) GenId() {
	e.Id = utils.ID("event")
}
