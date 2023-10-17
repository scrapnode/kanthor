package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Attempt struct {
	Publisher  streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Trigger    AttemptTrigger             `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
}

func (conf *Attempt) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.attempt.publisher: %v", err)
	}
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.attempt.subscriber: %v", err)
	}
	if err := conf.Trigger.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger: %v", err)
	}

	return nil
}

type AttemptTrigger struct {
	Plan AttemptTriggerPlan `json:"plan" yaml:"plan" mapstructure:"plan"`
	Exec AttemptTriggerExec `json:"exec" yaml:"exec" mapstructure:"exec"`
}

func (conf *AttemptTrigger) Validate() error {
	if err := conf.Plan.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.plan: %v", err)
	}
	if err := conf.Exec.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.exec: %v", err)
	}
	return nil
}

type AttemptTriggerPlan struct {
	LockDuration int64 `json:"lock_duration" yaml:"lock_duration" mapstructure:"lock_duration"`
	Timeout      int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	RateLimit    int   `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_to" yaml:"scan_to" mapstructure:"scan_to"`
}

func (conf *AttemptTriggerPlan) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.lock_duration", conf.LockDuration, 15000),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.timeout", conf.Timeout, 3000),
		validator.NumberGreaterThan("config.attempt.trigger.cron.rate_limit", conf.RateLimit, 0),
		validator.NumberGreaterThan("scan_start", conf.ScanStart, conf.ScanEnd),
		validator.NumberLessThan("scan_to", conf.ScanEnd, conf.ScanStart),
	)
}

type AttemptTriggerExec struct {
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	RateLimit int   `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
	Delay     int64 `json:"delay" yaml:"delay" mapstructure:"delay"`
}

func (conf *AttemptTriggerExec) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.exec.timeout", conf.Timeout, 3000),
		validator.NumberGreaterThan("config.attempt.trigger.exec.rate_limit", conf.RateLimit, 0),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.exec.delay", conf.Delay, 3000),
	)
}
