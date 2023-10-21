package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Dispatcher struct {
	Forwarder DispatcherForwarder `json:"forwarder" yaml:"forwarder" mapstructure:"forwarder"`
}

func (conf *Dispatcher) Validate() error {
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
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *DispatcherForwarderSend) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.dispatcher.forwarder.send.concurrency", conf.Concurrency, 0),
	)
}
