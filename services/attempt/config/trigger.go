package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Trigger struct {
	Planner  TriggerPlanner  `json:"planner" yaml:"planner" mapstructure:"planner"`
	Executor TriggerExecutor `json:"executor" yaml:"executor" mapstructure:"executor"`
}

func (conf *Trigger) Validate() error {
	if err := conf.Planner.Validate(); err != nil {
		return fmt.Errorf("attempt.trigger.planner: %v", err)
	}
	if err := conf.Executor.Validate(); err != nil {
		return fmt.Errorf("attempt.trigger.executor: %v", err)
	}
	return nil
}

type TriggerPlanner struct {
	Schedule string `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Timeout  int64  `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Size     int    `json:"size" yaml:"size" mapstructure:"size"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_end" yaml:"scan_end" mapstructure:"scan_end"`
}

func (conf *TriggerPlanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("attempt.trigger.planner.schedule", conf.Schedule),
		validator.NumberGreaterThanOrEqual("attempt.trigger.planner.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("attempt.trigger.executor.size", conf.Size, 0),
		validator.NumberLessThan("attempt.trigger.planner.scan_end", conf.ScanEnd, 0),
		validator.NumberLessThan("attempt.trigger.planner.scan_start", conf.ScanStart, conf.ScanEnd),
	)
}

type TriggerExecutor struct {
	Timeout      int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Concurrency  int   `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
	ArrangeDelay int64 `json:"arrange_delay" yaml:"arrange_delay" mapstructure:"arrange_delay"`
}

func (conf *TriggerExecutor) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("attempt.trigger.executor.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("attempt.trigger.executor.concurrency", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("attempt.trigger.executor.arrange_delay", conf.ArrangeDelay, 60000),
	)
}
