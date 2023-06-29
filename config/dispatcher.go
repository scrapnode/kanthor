package config

import "github.com/scrapnode/kanthor/infrastructure/streaming"

type Dispatcher struct {
	Consumer streaming.SubscriberConfig `json:"consumer" mapstructure:"consumer"`
	Metrics  Server                     `json:"metrics" mapstructure:"metrics"`
}
