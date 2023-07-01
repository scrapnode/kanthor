package config

import (
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Scheduler struct {
	Publisher  streaming.PublisherConfig  `json:"publisher" mapstructure:"publisher"`
	Subscriber streaming.SubscriberConfig `json:"subscriber" mapstructure:"subscriber"`
	Cache      *cache.Config              `json:"cache" mapstructure:"cache"`
}
