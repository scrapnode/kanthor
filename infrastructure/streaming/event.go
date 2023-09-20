package streaming

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/scrapnode/kanthor/pkg/validator"
)

var (
	MetaAppId = "KANTHOR_META_APP_ID"
	MetaType  = "KANTHOR_META_TYPE"
	MetaId    = "KANTHOR_META_ID"
)

type Event struct {
	Subject string `json:"subject"`
	AppId   string `json:"app_id"`
	Type    string `json:"type"`

	Id       string            `json:"id"`
	Data     []byte            `json:"data"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Event) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("subject", e.Subject),
		validator.StringRequired("app_id", e.AppId),
		validator.StringRequired("type", e.Type),
		validator.StringRequired("id", e.Id),
		validator.SliceRequired("data", e.Data),
		validator.MapNotNil("metadata", e.Metadata),
	)
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

func (e *Event) Is(topic string) bool {
	r := regexp.MustCompile(fmt.Sprintf("%s.([a-zA-z0-9]+).%s", Namespace, topic))
	return r.MatchString(e.Subject)
}

func Subject(ns string, tier string, topic string, segments ...string) string {
	return strings.Join(append([]string{ns, tier, topic}, segments...), ".")
}
