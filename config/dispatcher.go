package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

type Dispatcher struct {
	Publisher      streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber     streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Sender         sender.Config              `json:"sender" yaml:"sender" mapstructure:"sender"`
	Cache          cache.Config               `json:"cache" yaml:"cache" mapstructure:"cache"`
	CircuitBreaker circuitbreaker.Config      `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
	Metrics        metric.Config              `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
}

func (conf *Dispatcher) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.publisher: %v", err)
	}
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.subscriber: %v", err)
	}
	if err := conf.Sender.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.sender: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.cache: %v", err)
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.circuit_breaker: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.metrics: %v", err)
	}

	return nil
}
