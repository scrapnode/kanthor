package config

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Endeavor struct {
	Planner  EndeavorPlanner  `json:"planner" yaml:"planner" mapstructure:"planner"`
	Executor EndeavorExecutor `json:"executor" yaml:"executor" mapstructure:"executor"`
}

func (conf *Endeavor) Validate() error {
	if err := conf.Planner.Validate(); err != nil {
		return err
	}
	if err := conf.Executor.Validate(); err != nil {
		return err
	}
	return nil
}

type EndeavorPlanner struct {
	Schedule string `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Timeout  int64  `json:"timeout" yaml:"timeout" mapstructure:"timeout"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_end" yaml:"scan_end" mapstructure:"scan_end"`
}

func (conf *EndeavorPlanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("CONFIG.ATTEMPT.ENDEAVOR.PLANNER.SCHEDULE", conf.Schedule),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.ENDEAVOR.PLANNER.TIMEOUT", conf.Timeout, 1000),
		validator.NumberLessThan("CONFIG.ATTEMPT.ENDEAVOR.PLANNER.SCAN_END", conf.ScanEnd, 0),
		validator.NumberLessThan("CONFIG.ATTEMPT.ENDEAVOR.PLANNER.SCAN_START", conf.ScanStart, conf.ScanEnd),
	)
}

type EndeavorExecutor struct {
	Timeout         int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Concurrency     int   `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
	RescheduleDelay int64 `json:"reschedule_delay" yaml:"reschedule_delay" mapstructure:"reschedule_delay"`
}

func (conf *EndeavorExecutor) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.ENDEAVOR.EXECUTOR.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.ENDEAVOR.EXECUTOR.CONCURRENCY", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.ENDEAVOR.EXECUTOR.RESCHEDULE_DELAY", conf.RescheduleDelay, 60000),
	)
}
