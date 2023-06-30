package config

import (
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

type Dispatcher struct {
	Consumer streaming.SubscriberConfig `json:"consumer" mapstructure:"consumer"`
	Sender   sender.Config              `json:"sender" mapstructure:"sender"`
	Metrics  Server                     `json:"metrics" mapstructure:"metrics"`
}
