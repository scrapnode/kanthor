package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Scheduler struct {
	Publisher  streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Cache      cache.Config               `json:"cache" yaml:"cache"`
	Metrics    metric.Config              `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
}

func (conf *Scheduler) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.publisher: %v", err)
	}
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.subscriber: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.metrics: %v", err)
	}

	return nil
}
