package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Dispatcher, provider.Unmarshal(&conf)
}

type Wrapper struct {
	Dispatcher Config `json:"dispatcher" yaml:"dispatcher" mapstructure:"dispatcher"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Dispatcher.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Forwarder DispatcherForwarder `json:"forwarder" yaml:"forwarder" mapstructure:"forwarder"`
}

func (conf *Config) Validate() error {
	if err := conf.Forwarder.Validate(); err != nil {
		return err
	}

	return nil
}

type DispatcherForwarder struct {
	Send DispatcherForwarderSend `json:"send" yaml:"send" mapstructure:"send"`
}

func (conf *DispatcherForwarder) Validate() error {
	if err := conf.Send.Validate(); err != nil {
		return err
	}
	return nil
}

type DispatcherForwarderSend struct {
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *DispatcherForwarderSend) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("CONFIG.DISPATCHER.FORWARDER.SEND.CONCURRENCY", conf.Concurrency, 0),
	)
}
