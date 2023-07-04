package config

import (
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Dataplane struct {
	GRPC      Server                    `json:"grpc" mapstructure:"grpc"`
	Publisher streaming.PublisherConfig `json:"publisher" mapstructure:"publisher"`
	Cache     *cache.Config             `json:"cache" mapstructure:"cache"`

	Metrics metric.Config `json:"metrics" mapstructure:"metrics"`
}
