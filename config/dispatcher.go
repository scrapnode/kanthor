package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

type Dispatcher struct {
	Publisher      streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher" validate:"required"`
	Subscriber     streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber" validate:"required"`
	Sender         sender.Config              `json:"sender" yaml:"sender" mapstructure:"sender" validate:"required"`
	Cache          cache.Config               `json:"cache" yaml:"cache" mapstructure:"cache" validate:"required"`
	CircuitBreaker circuitbreaker.Config      `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker" validate:"required"`
}

func (conf *Dispatcher) Validate() error {
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
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Cache: %v", err)
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return fmt.Errorf("config.Dispatcher.CircuitBreaker: %v", err)
	}

	return nil
}
