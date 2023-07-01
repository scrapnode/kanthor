package config

import (
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Dataplane struct {
	GRPC      Server                    `json:"grpc" mapstructure:"grpc"`
	Publisher streaming.PublisherConfig `json:"publisher" mapstructure:"publisher"`
	Cache     *cache.Config             `json:"cache" mapstructure:"cache"`
	Metrics   Server                    `json:"metrics" mapstructure:"metrics"`
}
