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
		return fmt.Errorf("config.Storage.Subscriber: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Storage.Subscriber: %v", err)
	}

	return nil
}
