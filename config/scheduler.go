package config

import "github.com/scrapnode/kanthor/infrastructure/streaming"

type Scheduler struct {
	Consumer streaming.SubscriberConfig `json:"consumer" mapstructure:"consumer"`
}
