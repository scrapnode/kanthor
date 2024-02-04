package streaming

import (
	"encoding/json"

	"github.com/scrapnode/kanthor/pkg/validator"
)

var (
	MetaId               = "KANTHOR_META_ID"
	HeaderTelemetryTrace = "x-telemtry-trace"
)

type Event struct {
	Subject string `json:"subject"`

	Id       string            `json:"id"`
	Data     []byte            `json:"data"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Event) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("subject", e.Subject),
		validator.StringRequired("id", e.Id),
		validator.SliceRequired("data", e.Data),
		validator.MapNotNil("metadata", e.Metadata),
	)
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}
