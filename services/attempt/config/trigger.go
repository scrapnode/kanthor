package config

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Trigger struct {
	Planner  TriggerPlanner  `json:"planner" yaml:"planner" mapstructure:"planner"`
	Executor TriggerExecutor `json:"executor" yaml:"executor" mapstructure:"executor"`
}

func (conf *Trigger) Validate() error {
	if err := conf.Planner.Validate(); err != nil {
		return err
	}
	if err := conf.Executor.Validate(); err != nil {
		return err
	}
	return nil
}

type TriggerPlanner struct {
	Schedule string `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Size     int    `json:"size" yaml:"size" mapstructure:"size"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_end" yaml:"scan_end" mapstructure:"scan_end"`
}

func (conf *TriggerPlanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("CONFIG.ATTEMPT.TRIGGER.PLANNER.SCHEDULE", conf.Schedule),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.TRIGGER.PLANNER.SIZE", conf.Size, 0),
		validator.NumberLessThan("CONFIG.ATTEMPT.TRIGGER.PLANNER.SCAN_END", conf.ScanEnd, 0),
		validator.NumberLessThan("CONFIG.ATTEMPT.TRIGGER.PLANNER.SCAN_START", conf.ScanStart, conf.ScanEnd),
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
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.TRIGGER.EXECUTOR.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.TRIGGER.EXECUTOR.CONCURRENCY", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.TRIGGER.EXECUTOR.ARRANGE_DELAY", conf.ArrangeDelay, 60000),
	)
}
