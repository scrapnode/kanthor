package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Dispatcher struct {
	Publisher      streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber     streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Sender         sender.Config              `json:"sender" yaml:"sender" mapstructure:"sender"`
	Cache          cache.Config               `json:"cache" yaml:"cache" mapstructure:"cache"`
	CircuitBreaker circuitbreaker.Config      `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
	Metrics        metric.Config              `json:"metrics" yaml:"metrics" mapstructure:"metrics"`

	Forwarder DispatcherForwarder `json:"forwarder" yaml:"forwarder" mapstructure:"forwarder"`
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
	if err := conf.Forwarder.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.forwarder: %v", err)
	}

	return nil
}

type DispatcherForwarder struct {
	Send DispatcherForwarderSend `json:"send" yaml:"send" mapstructure:"send"`
}

func (conf *DispatcherForwarder) Validate() error {
	if err := conf.Send.Validate(); err != nil {
		return fmt.Errorf("config.dispatcher.forwarder.send: %v", err)
	}
	return nil
}

type DispatcherForwarderSend struct {
	ChunkSize    int   `json:"chunk_size" yaml:"chunk_size" mapstructure:"chunk_size"`
	ChunkTimeout int64 `json:"chunk_timeout" yaml:"chunk_timeout" mapstructure:"chunk_timeout"`
}

func (conf *DispatcherForwarderSend) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.dispatcher.forwarder.send.chunk_size", conf.ChunkSize, 0),
		validator.NumberGreaterThanOrEqual("config.dispatcher.forwarder.send.chunk_timeout", conf.ChunkTimeout, 1000),
	)
}
