package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Attempt struct {
	Trigger AttemptTrigger `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
}

func (conf *Attempt) Validate() error {
	if err := conf.Trigger.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger: %v", err)
	}

	return nil
}

type AttemptTrigger struct {
	Publisher streaming.PublisherConfig `json:"publisher" yaml:"publisher" mapstructure:"publisher"`

	Planner  AttemptTriggerPlanner  `json:"planer" yaml:"planer" mapstructure:"planer"`
	Executor AttemptTriggerExecutor `json:"executor" yaml:"executor" mapstructure:"executor"`
}

func (conf *AttemptTrigger) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.publisher: %v", err)
	}
	if err := conf.Planner.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.planner: %v", err)
	}
	if err := conf.Executor.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.executor: %v", err)
	}
	return nil
}

type AttemptTriggerPlanner struct {
	Schedule  string `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Timeout   int64  `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	RateLimit int    `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_end" yaml:"scan_end" mapstructure:"scan_end"`
}

func (conf *AttemptTriggerPlanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("config.attempt.trigger.planner.schedule", conf.Schedule),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.planner.timeout", conf.Timeout, 3000),
		validator.NumberGreaterThan("config.attempt.trigger.planner.rate_limit", conf.RateLimit, 0),
		validator.NumberGreaterThan("config.attempt.trigger.planner.scan_start", conf.ScanStart, conf.ScanEnd),
		validator.NumberLessThan("config.attempt.trigger.planner.scan_end", conf.ScanEnd, conf.ScanStart),
	)
}

type AttemptTriggerExecutor struct {
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`

	Timeout      int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	RateLimit    int   `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
	AttemptDelay int64 `json:"attempt_delay" yaml:"attempt_delay" mapstructure:"attempt_delay"`
}

func (conf *AttemptTriggerExecutor) Validate() error {
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.executor.subscriber: %v", err)
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.executor.timeout", conf.Timeout, 3000),
		validator.NumberGreaterThan("config.attempt.trigger.executor.rate_limit", conf.RateLimit, 0),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.executor.attempt_delay", conf.AttemptDelay, 1000),
	)
}
