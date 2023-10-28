package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Attempt, provider.Unmarshal(&conf)
}

type Wrapper struct {
	Attempt Config `json:"attempt" yaml:"attempt" mapstructure:"attempt"`
}

type Config struct {
	Trigger AttemptTrigger `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
}

func (conf *Config) Validate() error {
	if err := conf.Trigger.Validate(); err != nil {
		return fmt.Errorf("attempt.trigger: %v", err)
	}

	return nil
}

type AttemptTrigger struct {
	Planner  AttemptTriggerPlanner  `json:"planner" yaml:"planner" mapstructure:"planner"`
	Executor AttemptTriggerExecutor `json:"executor" yaml:"executor" mapstructure:"executor"`
}

func (conf *AttemptTrigger) Validate() error {
	if err := conf.Planner.Validate(); err != nil {
		return fmt.Errorf("attempt.trigger.planner: %v", err)
	}
	if err := conf.Executor.Validate(); err != nil {
		return fmt.Errorf("attempt.trigger.executor: %v", err)
	}
	return nil
}

type AttemptTriggerPlanner struct {
	Schedule string `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Timeout  int64  `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Size     int    `json:"size" yaml:"size" mapstructure:"size"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_end" yaml:"scan_end" mapstructure:"scan_end"`
}

func (conf *AttemptTriggerPlanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("attempt.trigger.planner.schedule", conf.Schedule),
		validator.NumberGreaterThanOrEqual("attempt.trigger.planner.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("attempt.trigger.planner.size", conf.Size, 0),
		validator.NumberLessThan("attempt.trigger.planner.scan_end", conf.ScanEnd, 0),
		validator.NumberLessThan("attempt.trigger.planner.scan_start", conf.ScanStart, conf.ScanEnd),
	)
}

type AttemptTriggerExecutor struct {
	Timeout      int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Size         int   `json:"size" yaml:"size" mapstructure:"size"`
	AttemptDelay int64 `json:"attempt_delay" yaml:"attempt_delay" mapstructure:"attempt_delay"`
}

func (conf *AttemptTriggerExecutor) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("attempt.trigger.executor.timeout", int(conf.Timeout), 1000),
		validator.NumberGreaterThan("attempt.trigger.executor.size", conf.Size, 0),
		validator.NumberGreaterThanOrEqual("attempt.trigger.executor.attempt_delay", conf.AttemptDelay, 60000),
	)
}
