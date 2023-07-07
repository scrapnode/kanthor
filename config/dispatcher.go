package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

type Dispatcher struct {
	Publisher      streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher" validate:"required"`
	Subscriber     streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber" validate:"required"`
	Sender         sender.Config              `json:"sender" yaml:"sender" mapstructure:"sender" validate:"required"`
	Cache          *cache.Config              `json:"cache" yaml:"cache" mapstructure:"cache" validate:"-"`
	CircuitBreaker circuitbreaker.Config      `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker" validate:"required"`

	Metrics metric.Config `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"-"`
}

func (conf Dispatcher) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Dispatcher: %v", err)
	}

	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.Dispatcher.Publisher: %v", err)
	}
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.Dispatcher.Subscriber: %v", err)
	}
	if err := conf.Sender.Validate(); err != nil {
		return fmt.Errorf("config.Dispatcher.Sender: %v", err)
	}
	if conf.Cache != nil {
		if err := conf.Cache.Validate(); err != nil {
			return fmt.Errorf("config.Dispatcher.Cache: %v", err)
		}
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return fmt.Errorf("config.Dispatcher.CircuitBreaker: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Dispatcher.Metrics: %v", err)
	}

	return nil
}
