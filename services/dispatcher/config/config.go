package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Dispatcher Dispatcher `json:"dispatcher" yaml:"dispatcher" mapstructure:"dispatcher"`
}

func (conf *Config) Validate() error {
	if err := conf.Dispatcher.Validate(); err != nil {
		return err
	}
	return nil
}

type Dispatcher struct {
	Forwarder DispatcherForwarder `json:"forwarder" yaml:"forwarder" mapstructure:"forwarder"`
}

func (conf *Dispatcher) Validate() error {
	if err := conf.Forwarder.Validate(); err != nil {
		return fmt.Errorf("dispatcher.forwarder: %v", err)
	}

	return nil
}

type DispatcherForwarder struct {
	Send DispatcherForwarderSend `json:"send" yaml:"send" mapstructure:"send"`
}

func (conf *DispatcherForwarder) Validate() error {
	if err := conf.Send.Validate(); err != nil {
		return fmt.Errorf("dispatcher.forwarder.send: %v", err)
	}
	return nil
}

type DispatcherForwarderSend struct {
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *DispatcherForwarderSend) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("dispatcher.forwarder.send.concurrency", conf.Concurrency, 0),
	)
}
