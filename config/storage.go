package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Storage struct {
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Metrics    metric.Config              `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
}

func (conf *Storage) Validate() error {
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.storage.subscriber: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.storage.metrics: %v", err)
	}

	return nil
}
