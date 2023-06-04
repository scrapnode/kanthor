package msgbroker

import (
	"context"
	"encoding/json"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (MsgBroker, error) {
	return NewNats(conf, logger)
}

type MsgBroker interface {
	patterns.Connectable
	Pub(ctx context.Context, event *Event) error
	Sub(ctx context.Context, handler Handler) error
}

type Event struct {
	Tier  string `json:"tier"`
	AppId string `json:"app_id"`
	Type  string `json:"type"`

	Id       string            `json:"id"`
	Data     []byte            `json:"data"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

type Handler func(event *Event) error

var (
	MetaTier  = "KANTHOR_META_TIER"
	MetaAppId = "KANTHOR_META_APP_ID"
	MetaType  = "KANTHOR_META_TYPE"
	MetaId    = "KANTHOR_META_ID"
)
