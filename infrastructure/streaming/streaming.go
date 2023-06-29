package streaming

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/infrastructure/utils"
	"strings"
)

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
	if e.Id == "" {
		e.Id = utils.ID("event")
	}
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

func Subject(segments ...string) string {
	return strings.Join(segments, ".")
}
