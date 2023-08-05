package streaming

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	MetaAppId = "KANTHOR_META_APP_ID"
	MetaType  = "KANTHOR_META_TYPE"
	MetaId    = "KANTHOR_META_ID"
)

type Event struct {
	Subject string `json:"subject" validate:"required"`
	AppId   string `json:"app_id" validate:"required"`
	Type    string `json:"type" validate:"required"`

	Id       string            `json:"id" validate:"required"`
	Data     []byte            `json:"data" validate:"required"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Event) Validate() error {
	return validator.New().Struct(e)
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

func Subject(ns string, tier string, topic string, segments ...string) string {
	return strings.Join(append([]string{ns, tier, topic}, segments...), ".")
}
