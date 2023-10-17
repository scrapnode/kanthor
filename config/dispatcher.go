package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Dispatcher struct {
	Publisher  streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Sender     sender.Config              `json:"sender" yaml:"sender" mapstructure:"sender"`

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
	RateLimit int   `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (conf *DispatcherForwarderSend) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.dispatcher.forwarder.send.rate_limit", conf.RateLimit, 0),
		validator.NumberGreaterThanOrEqual("config.dispatcher.forwarder.send.timeout", conf.Timeout, 1000),
	)
}
